// ABOUTME: Contract tests for the Spacedock-owned gate-binding boundary.
// ABOUTME: They validate the schema, dogfood fixtures, and additive-field preservation without runtime code.
package contract

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
	"unicode/utf16"
	"unicode/utf8"
)

const gateBindingSchemaPath = "../../docs/schema/gate-binding.v1.schema.json"
const helmSchemaSourceRevision = "3d39fcdc67cc6aac22d51f4abcc2dfdadf56c838"

func TestGateBindingSchemaPinsProviderTargetIdentity(t *testing.T) {
	schema := readJSONObject(t, gateBindingSchemaPath)

	assertStringValue(t, schema, "$schema", "https://json-schema.org/draft/2020-12/schema")
	assertStringValue(t, schema, "$id", "https://schemas.spacedock.dev/gate-binding.v1.schema.json")
	assertStringValue(t, schema, "type", "object")
	if got := schema["additionalProperties"]; got != true {
		t.Fatalf("root additionalProperties = %#v, want true for additive compatibility", got)
	}

	required := stringSet(t, schema["required"])
	for _, name := range []string{
		"ns",
		"entity_ref",
	} {
		if !required[name] {
			t.Errorf("schema required does not contain %q", name)
		}
	}
	for _, temporal := range []string{"gate_id", "resolution_id", "application_id", "application"} {
		if required[temporal] {
			t.Errorf("binding requires later lifecycle field %q", temporal)
		}
	}

	properties := objectValue(t, schema, "properties")
	assertStringValue(t, objectValue(t, properties, "ns"), "const", "spacedock.gate-binding.v1")
	for _, routed := range []string{"entity_ref", "stage", "target_stage", "workflow_ref", "provider_instance_id", "expected_revision"} {
		property := objectValue(t, properties, routed)
		if _, ok := property["pattern"].(string); !ok {
			t.Errorf("routing field %q lacks non-whitespace pattern", routed)
		}
	}
	for _, optional := range []string{"stage", "target_stage"} {
		if required[optional] {
			t.Errorf("existing two-field binding compatibility requires %q to stay optional", optional)
		}
	}
	for _, temporal := range []string{"gate_id", "resolution_id", "application_id", "application"} {
		if _, defined := properties[temporal]; defined {
			t.Errorf("binding schema defines later lifecycle field %q", temporal)
		}
	}
}

func TestGateBindingRejectsWhitespaceRoutingFields(t *testing.T) {
	schema := readJSONObject(t, gateBindingSchemaPath)
	for _, fixture := range []string{
		"invalid-whitespace-entity-ref.json",
		"invalid-whitespace-stage.json",
		"invalid-whitespace-target-stage.json",
		"invalid-whitespace-workflow-ref.json",
		"invalid-whitespace-provider-instance-id.json",
		"invalid-whitespace-expected-revision.json",
	} {
		t.Run(strings.TrimSuffix(fixture, ".json"), func(t *testing.T) {
			binding := readJSONObject(t, filepath.Join("testdata", "ledger-boundary", fixture))
			if err := validateGateBinding(schema, binding); err == nil {
				t.Fatal("whitespace-only routing field unexpectedly validates")
			}
		})
	}
}

func TestGateBindingDigestContractPinsRFC8785(t *testing.T) {
	rawSpec, err := os.ReadFile("../../docs/specs/ledger-gate-binding.md")
	if err != nil {
		t.Fatalf("read binding spec: %v", err)
	}
	for _, required := range []string{"RFC 8785", "JSON Canonicalization Scheme (JCS)"} {
		if !strings.Contains(string(rawSpec), required) {
			t.Errorf("binding digest contract does not name %q", required)
		}
	}

	unicodeVector := map[string]any{
		"€":      "Euro Sign",
		"\r":     "Carriage Return",
		"דּ":      "Hebrew Letter Dalet With Dagesh",
		"1":      "One",
		"😀":      "Emoji: Grinning Face",
		"\u0080": "Control",
		"ö":      "Latin Small Letter O With Diaeresis",
	}
	wantUnicode := "{\"\\r\":\"Carriage Return\",\"1\":\"One\",\"\u0080\":\"Control\",\"ö\":\"Latin Small Letter O With Diaeresis\",\"€\":\"Euro Sign\",\"😀\":\"Emoji: Grinning Face\",\"דּ\":\"Hebrew Letter Dalet With Dagesh\"}"
	assertJCSVector(t, "unicode-key-order-and-escaping", unicodeVector, wantUnicode)

	numberVector := map[string]any{
		"integer":      1.0,
		"large":        1e20,
		"negativeZero": math.Copysign(0, -1),
		"scientific":   1e21,
		"small":        1e-7,
		"threshold":    1e-6,
		"escape":       "\b\t\n\f\r\"\\",
	}
	wantNumber := "{\"escape\":\"\\b\\t\\n\\f\\r\\\"\\\\\",\"integer\":1,\"large\":100000000000000000000,\"negativeZero\":0,\"scientific\":1e+21,\"small\":1e-7,\"threshold\":0.000001}"
	assertJCSVector(t, "number-and-string-representation", numberVector, wantNumber)
}

func TestGateBindingHelmWireFixturesConformToPinnedSchemas(t *testing.T) {
	registry := loadPinnedHelmSchemas(t)
	root := filepath.Join("testdata", "ledger-boundary")

	recovery := readJSONObject(t, filepath.Join(root, "application-recovery.json"))
	validatePinnedHelmObject(t, registry, "https://helm.spacedock.dev/schemas/ledger/application-record-committed-request.v1.json", objectValue(t, recovery, "committed_request"))
	validatePinnedHelmObject(t, registry, "https://helm.spacedock.dev/schemas/ledger/application-committed-receipt.v1.json", objectValue(t, recovery, "accepted_receipt"))
	validatePinnedHelmObject(t, registry, "https://helm.spacedock.dev/schemas/ledger/application-committed-receipt.v1.json", objectValue(t, recovery, "same_key_same_body_replay"))
	conflict := objectValue(t, recovery, "changed_body_conflict")
	validatePinnedHelmObject(t, registry, "https://helm.spacedock.dev/schemas/ledger/application-record-committed-request.v1.json", objectValue(t, conflict, "request"))
	recoveryView := objectValue(t, objectValue(t, recovery, "response_lost_recovery"), "application_view")
	validatePinnedHelmObject(t, registry, "https://helm.spacedock.dev/schemas/ledger/application-view.v1.json", recoveryView)

	linkage := readJSONObject(t, filepath.Join(root, "valid-application-linkage.json"))
	validatePinnedHelmObject(t, registry, "https://helm.spacedock.dev/schemas/ledger/application-committed.v1.json", objectValue(t, linkage, "committed"))
	validatePinnedHelmObject(t, registry, "https://helm.spacedock.dev/schemas/ledger/application-observed.v1.json", objectValue(t, linkage, "observed"))
	validatePinnedHelmObject(t, registry, "https://helm.spacedock.dev/schemas/ledger/application-source-superseded.v1.json", objectValue(t, linkage, "source_superseded"))

	projection := readJSONObject(t, filepath.Join(root, "projection-recovery.json"))
	validatePinnedHelmObject(t, registry, "https://helm.spacedock.dev/schemas/ledger/projection-rewrite-quarantined.v1.json", objectValue(t, projection, "rewrite_quarantined"))
}

func TestGateBindingAcceptsExistingMinimalOpaqueSlot(t *testing.T) {
	schema := readJSONObject(t, gateBindingSchemaPath)
	binding := map[string]any{
		"ns":         "spacedock.gate-binding.v1",
		"entity_ref": "docs/ship-flow/m4.13",
	}
	if err := validateGateBinding(schema, binding); err != nil {
		t.Fatalf("existing Helm opaque binding slot rejected: %v", err)
	}
}

func TestGateBindingApplicationReplayAndResponseLostRecovery(t *testing.T) {
	fixture := readJSONObject(t, filepath.Join("testdata", "ledger-boundary", "application-recovery.json"))
	request := objectValue(t, fixture, "committed_request")
	assertDigestMatchesBinding(t, request, "fact", "body_digest")
	receipt := objectValue(t, fixture, "accepted_receipt")
	replay := objectValue(t, fixture, "same_key_same_body_replay")
	if !reflect.DeepEqual(receipt, replay) {
		t.Fatalf("same-key/same-body replay is not byte-semantic receipt replay\nreceipt: %#v\nreplay: %#v", receipt, replay)
	}
	assertEqualStringFields(t, request, receipt, "idempotency_key")
	assertEqualStringFields(t, request, receipt, "body_digest")

	conflict := objectValue(t, fixture, "changed_body_conflict")
	conflictingRequest := objectValue(t, conflict, "request")
	assertDigestMatchesBinding(t, conflictingRequest, "fact", "body_digest")
	assertEqualStringFields(t, request, conflictingRequest, "idempotency_key")
	if request["body_digest"] == conflictingRequest["body_digest"] {
		t.Fatal("changed-body conflict fixture must change body_digest")
	}
	conflictError := objectValue(t, conflict, "error")
	assertNumberValue(t, conflictError, "status", 409)
	assertStringValue(t, conflictError, "code", "helm.application.idempotency_conflict.v1")
	assertStringValue(t, conflict, "effect", "refuse_without_new_fact")

	recovery := objectValue(t, fixture, "response_lost_recovery")
	assertStringValue(t, recovery, "coordinator_state", "commit_succeeded_attestation_outcome_unknown")
	read := objectValue(t, recovery, "read_application")
	view := objectValue(t, recovery, "application_view")
	assertEqualStringFields(t, read, view, "application_id")
	assertStringValue(t, view, "schema", "helm.application.view.v1")
	assertStringValue(t, view, "status", "applied")
	viewReceipt := objectValue(t, view, "committed_receipt")
	if !reflect.DeepEqual(receipt, viewReceipt) {
		t.Fatal("response-lost recovery view must return the accepted committed receipt")
	}
}

func TestGateBindingProjectionEpochInvalidationAndRewriteFreeze(t *testing.T) {
	fixture := readJSONObject(t, filepath.Join("testdata", "ledger-boundary", "projection-recovery.json"))
	page := objectValue(t, fixture, "initial_page")
	invalidation := objectValue(t, fixture, "cursor_invalidation")
	request := objectValue(t, invalidation, "request")
	assertEqualStringFields(t, page, request, "projection_epoch")
	pageCursor, pageCursorOK := page["next_cursor"].(string)
	requestAfter, requestAfterOK := request["after"].(string)
	if !pageCursorOK || !requestAfterOK || pageCursor != requestAfter {
		t.Fatalf("event poll must replay opaque next_cursor as after: next_cursor=%#v after=%#v", page["next_cursor"], request["after"])
	}
	errorObject := objectValue(t, invalidation, "error")
	assertNumberValue(t, errorObject, "status", 409)
	assertStringValue(t, errorObject, "code", "helm.projection.cursor_invalidated.v1")
	details := objectValue(t, errorObject, "details")
	if details["current_projection_epoch"] == page["projection_epoch"] {
		t.Fatal("cursor invalidation must return a new current projection epoch")
	}

	exposed := objectValue(t, fixture, "exposed_binding")
	candidate := objectValue(t, fixture, "candidate_binding")
	quarantine := objectValue(t, fixture, "rewrite_quarantined")
	assertStringValue(t, quarantine, "schema", "helm.projection.rewrite_quarantined.v1")
	assertStringValue(t, quarantine, "old_projection_epoch", page["projection_epoch"].(string))
	assertStringValue(t, quarantine, "new_projection_epoch", details["current_projection_epoch"].(string))
	assertStringValue(t, quarantine, "invalidated_cursor", pageCursor)
	assertDigestMatchesBinding(t, quarantine, "exposed_binding", "exposed_binding_digest")
	assertDigestMatchesBinding(t, quarantine, "candidate_binding", "candidate_binding_digest")
	if !reflect.DeepEqual(exposed, objectValue(t, quarantine, "exposed_binding")) {
		t.Fatal("rewrite quarantine must freeze the previously exposed binding")
	}
	if reflect.DeepEqual(exposed, candidate) {
		t.Fatal("rewrite fixture must select a different candidate binding")
	}
	if !reflect.DeepEqual(candidate, objectValue(t, quarantine, "candidate_binding")) {
		t.Fatal("rewrite quarantine must retain the candidate for explicit reconciliation")
	}
	expected := objectValue(t, fixture, "expected")
	assertBoolValue(t, expected, "exposed_binding_frozen", true)
	assertBoolValue(t, expected, "new_application_minted", false)
}

func TestGateBindingProviderPinFailsClosedWithTypedError(t *testing.T) {
	schema := readJSONObject(t, gateBindingSchemaPath)
	fixture := readJSONObject(t, filepath.Join("testdata", "ledger-boundary", "provider-pin-refusals.json"))
	expectedPin := objectValue(t, fixture, "expected_pin")
	for index, rawCase := range arrayValue(t, fixture, "cases") {
		caseObject, ok := rawCase.(map[string]any)
		if !ok {
			t.Fatalf("case %d is not an object", index)
		}
		name, _ := caseObject["name"].(string)
		t.Run(name, func(t *testing.T) {
			observedPin := objectValue(t, caseObject, "observed_pin")
			binding := objectValue(t, caseObject, "binding")
			switch name {
			case "unsupported-version":
				if expectedPin["ns"] == observedPin["ns"] {
					t.Fatal("unsupported-version case does not change ns")
				}
				if err := validateGateBinding(schema, binding); err == nil {
					t.Fatal("unsupported binding namespace unexpectedly validates")
				}
			case "digest-mismatch":
				if expectedPin["digest"] == observedPin["digest"] {
					t.Fatal("digest-mismatch case does not change digest")
				}
				if err := validateGateBinding(schema, binding); err != nil {
					t.Fatalf("digest mismatch must not reinterpret otherwise valid binding fields: %v", err)
				}
			default:
				t.Fatalf("unknown provider pin refusal case %q", name)
			}
			errorObject := objectValue(t, caseObject, "error")
			assertNumberValue(t, errorObject, "status", 409)
			assertStringValue(t, errorObject, "code", "helm.provider_binding.unsupported_version.v1")
			assertStringValue(t, caseObject, "effect", "refuse_before_mutation")
			assertBoolValue(t, caseObject, "fallback_guess_allowed", false)
		})
	}
}

func TestGateBindingFixtures(t *testing.T) {
	schema := readJSONObject(t, gateBindingSchemaPath)
	cases := []struct {
		name    string
		fixture string
		valid   bool
	}{
		{name: "minimal", fixture: "valid-minimal-binding.json", valid: true},
		{name: "forward-compatible", fixture: "valid-binding-with-extensions.json", valid: true},
		{name: "wrong-provider", fixture: "invalid-provider-namespace.json", valid: false},
		{name: "missing-entity", fixture: "invalid-missing-entity-ref.json", valid: false},
		{name: "empty-entity", fixture: "invalid-empty-entity-ref.json", valid: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fixture := readJSONObject(t, filepath.Join("testdata", "ledger-boundary", tc.fixture))
			err := validateGateBinding(schema, fixture)
			if tc.valid && err != nil {
				t.Fatalf("valid fixture rejected: %v", err)
			}
			if !tc.valid && err == nil {
				t.Fatal("invalid fixture accepted")
			}
		})
	}
}

func TestGateBindingUnknownFieldsSurviveRoundTrip(t *testing.T) {
	original := readJSONObject(t, filepath.Join("testdata", "ledger-boundary", "valid-binding-with-extensions.json"))
	raw, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal fixture: %v", err)
	}

	var roundTripped map[string]any
	if err := json.Unmarshal(raw, &roundTripped); err != nil {
		t.Fatalf("unmarshal fixture: %v", err)
	}
	if !reflect.DeepEqual(original, roundTripped) {
		t.Fatalf("unknown additive fields changed during round trip\noriginal: %#v\nround trip: %#v", original, roundTripped)
	}

	future := objectValue(t, roundTripped, "spacedock:future_binding_hint")
	if future["mode"] == nil || roundTripped["spacedock:provider_instance_id"] == nil {
		t.Fatal("forward-compatible fixture must exercise unknown scalar and object fields")
	}
}

func TestGateBindingBoundaryFixtureKeepsLaterApplicationFactsOutsideBinding(t *testing.T) {
	schema := readJSONObject(t, gateBindingSchemaPath)
	boundary := readJSONObject(t, filepath.Join("testdata", "ledger-boundary", "valid-application-linkage.json"))
	gate := objectValue(t, boundary, "gate")
	binding := objectValue(t, gate, "binding")
	if err := validateGateBinding(schema, binding); err != nil {
		t.Fatalf("boundary binding rejected: %v", err)
	}
	for _, temporal := range []string{"gate_id", "resolution_id", "application_id", "application"} {
		if _, exists := binding[temporal]; exists {
			t.Errorf("binding contains later lifecycle field %q", temporal)
		}
	}

	resolution := objectValue(t, boundary, "resolution")
	application := objectValue(t, boundary, "application")
	committed := objectValue(t, boundary, "committed")
	observed := objectValue(t, boundary, "observed")
	assertEqualStringFields(t, gate, resolution, "gate_id")
	assertEqualStringFields(t, gate, application, "gate_id")
	assertEqualStringFields(t, resolution, application, "resolution_id")
	assertEqualStringFields(t, application, committed, "application_id")
	assertEqualStringFields(t, application, observed, "application_id")
	assertStringValue(t, committed, "schema", "helm.application.committed.v1")
	assertStringValue(t, observed, "schema", "helm.application.observed.v1")
	expected := objectValue(t, boundary, "expected")
	assertStringValue(t, expected, "committed_effect", "pending_apply_to_applied")
	assertStringValue(t, expected, "observed_effect", "audit_only")
}

func validateGateBinding(schema, binding map[string]any) error {
	properties, ok := schema["properties"].(map[string]any)
	if !ok {
		return fmt.Errorf("schema properties are missing")
	}
	for name := range stringSetValue(schema["required"]) {
		if _, exists := binding[name]; !exists {
			return fmt.Errorf("required property %q is missing", name)
		}
	}
	for name, rawPropertySchema := range properties {
		value, exists := binding[name]
		if !exists {
			continue
		}
		propertySchema, ok := rawPropertySchema.(map[string]any)
		if !ok {
			return fmt.Errorf("property schema %q is not an object", name)
		}
		if err := validateSchemaValue(name, propertySchema, value); err != nil {
			return err
		}
	}
	return nil
}

func validateSchemaValue(path string, schema map[string]any, value any) error {
	return validateSchemaValueWithRegistry(path, schema, value, nil)
}

func validateSchemaValueWithRegistry(path string, schema map[string]any, value any, registry map[string]map[string]any) error {
	if ref, ok := schema["$ref"].(string); ok {
		resolved, exists := registry[ref]
		if !exists {
			return fmt.Errorf("%s references unpinned schema %q", path, ref)
		}
		return validateSchemaValueWithRegistry(path, resolved, value, registry)
	}
	if constant, ok := schema["const"]; ok && !reflect.DeepEqual(value, constant) {
		return fmt.Errorf("%s = %#v, want constant %#v", path, value, constant)
	}
	if enum, ok := schema["enum"].([]any); ok {
		found := false
		for _, candidate := range enum {
			found = found || reflect.DeepEqual(value, candidate)
		}
		if !found {
			return fmt.Errorf("%s = %#v, not in enum %#v", path, value, enum)
		}
	}

	switch schema["type"] {
	case "string":
		text, ok := value.(string)
		if !ok {
			return fmt.Errorf("%s must be a string", path)
		}
		if minLength, ok := schema["minLength"].(float64); ok && len(text) < int(minLength) {
			return fmt.Errorf("%s is shorter than %d", path, int(minLength))
		}
		if pattern, ok := schema["pattern"].(string); ok {
			matched, err := regexp.MatchString(pattern, text)
			if err != nil {
				return fmt.Errorf("%s has invalid schema pattern: %w", path, err)
			}
			if !matched {
				return fmt.Errorf("%s = %q does not match %q", path, text, pattern)
			}
		}
		if format, ok := schema["format"].(string); ok && format == "date-time" {
			if _, err := time.Parse(time.RFC3339, text); err != nil {
				return fmt.Errorf("%s = %q is not RFC3339 date-time", path, text)
			}
		}
	case "object":
		object, ok := value.(map[string]any)
		if !ok {
			return fmt.Errorf("%s must be an object", path)
		}
		_ = object
	case "array":
		items, ok := value.([]any)
		if !ok {
			return fmt.Errorf("%s must be an array", path)
		}
		if itemSchema, ok := schema["items"].(map[string]any); ok {
			for index, item := range items {
				if err := validateSchemaValueWithRegistry(fmt.Sprintf("%s[%d]", path, index), itemSchema, item, registry); err != nil {
					return err
				}
			}
		}
	case "integer":
		number, ok := value.(float64)
		if !ok || math.Trunc(number) != number {
			return fmt.Errorf("%s must be an integer", path)
		}
		if minimum, ok := schema["minimum"].(float64); ok && number < minimum {
			return fmt.Errorf("%s = %v, want >= %v", path, number, minimum)
		}
	}

	if object, ok := value.(map[string]any); ok {
		for name := range stringSetValue(schema["required"]) {
			if _, exists := object[name]; !exists {
				return fmt.Errorf("%s.%s is required", path, name)
			}
		}
		properties, _ := schema["properties"].(map[string]any)
		for name, rawPropertySchema := range properties {
			child, exists := object[name]
			if !exists {
				continue
			}
			propertySchema, ok := rawPropertySchema.(map[string]any)
			if !ok {
				return fmt.Errorf("%s.%s schema is not an object", path, name)
			}
			if err := validateSchemaValueWithRegistry(path+"."+name, propertySchema, child, registry); err != nil {
				return err
			}
		}
	}

	if branches, ok := schema["allOf"].([]any); ok {
		for index, rawBranch := range branches {
			branch, ok := rawBranch.(map[string]any)
			if !ok {
				return fmt.Errorf("%s.allOf[%d] is not an object", path, index)
			}
			if err := validateSchemaValueWithRegistry(path, branch, value, registry); err != nil {
				return err
			}
		}
	}
	if condition, ok := schema["if"].(map[string]any); ok {
		if validateSchemaValueWithRegistry(path, condition, value, registry) == nil {
			if thenSchema, ok := schema["then"].(map[string]any); ok {
				if err := validateSchemaValueWithRegistry(path, thenSchema, value, registry); err != nil {
					return err
				}
			}
		} else if elseSchema, ok := schema["else"].(map[string]any); ok {
			if err := validateSchemaValueWithRegistry(path, elseSchema, value, registry); err != nil {
				return err
			}
		}
	}
	if notSchema, ok := schema["not"].(map[string]any); ok {
		if validateSchemaValueWithRegistry(path, notSchema, value, registry) == nil {
			return fmt.Errorf("%s matches forbidden schema", path)
		}
	}
	if choices, ok := schema["anyOf"].([]any); ok {
		matched := false
		for _, rawChoice := range choices {
			choice, ok := rawChoice.(map[string]any)
			if ok && validateSchemaValueWithRegistry(path, choice, value, registry) == nil {
				matched = true
				break
			}
		}
		if !matched {
			return fmt.Errorf("%s does not match any allowed schema", path)
		}
	}
	return nil
}

func loadPinnedHelmSchemas(t *testing.T) map[string]map[string]any {
	t.Helper()
	root := filepath.Join("testdata", "ledger-boundary", "helm-schemas")
	manifest := readJSONObject(t, filepath.Join(root, "manifest.json"))
	provenanceStatus, _ := manifest["provenance_status"].(string)
	if provenanceStatus != "complete" {
		t.Fatalf("Helm schema provenance_status = %q, want complete", provenanceStatus)
	}
	assertStringValue(t, manifest, "source_repository", "https://github.com/spacedock-dev/helm")
	assertStringValue(t, manifest, "source_revision", helmSchemaSourceRevision)
	assertGitObjectID(t, manifest, "source_revision")
	pins := objectValue(t, manifest, "schemas")
	if len(pins) != 8 {
		t.Fatalf("pinned Helm schema closure contains %d schemas, want 8", len(pins))
	}
	registry := make(map[string]map[string]any, len(pins))
	for filename, rawPin := range pins {
		pin, ok := rawPin.(map[string]any)
		if !ok {
			t.Fatalf("schema pin for %s is not an object", filename)
		}
		assertStringValue(t, pin, "source_path", filepath.ToSlash(filepath.Join("src", "ledger", "contracts", "schemas", filename)))
		assertGitObjectID(t, pin, "git_blob")
		digest, ok := pin["raw_sha256"].(string)
		if !ok {
			t.Fatalf("raw_sha256 for %s is not a string", filename)
		}
		raw, err := os.ReadFile(filepath.Join(root, filename))
		if err != nil {
			t.Fatalf("read pinned Helm schema %s: %v", filename, err)
		}
		got := fmt.Sprintf("sha256:%x", sha256.Sum256(raw))
		if got != digest {
			t.Fatalf("pinned Helm schema %s digest = %s, want %s", filename, got, digest)
		}
		var schema map[string]any
		if err := json.Unmarshal(raw, &schema); err != nil {
			t.Fatalf("parse pinned Helm schema %s: %v", filename, err)
		}
		id, ok := schema["$id"].(string)
		if !ok {
			t.Fatalf("pinned Helm schema %s lacks $id", filename)
		}
		registry[id] = schema
	}
	for sourceID, schema := range registry {
		refs := make(map[string]bool)
		collectSchemaRefs(schema, refs)
		for ref := range refs {
			if strings.HasPrefix(ref, "https://helm.spacedock.dev/schemas/") {
				if _, exists := registry[ref]; !exists {
					t.Fatalf("pinned Helm schema closure is incomplete: %s references absent %s", sourceID, ref)
				}
			}
		}
	}
	return registry
}

func assertGitObjectID(t *testing.T, object map[string]any, name string) {
	t.Helper()
	value, ok := object[name].(string)
	matched, err := regexp.MatchString("^[0-9a-f]{40}$", value)
	if !ok || err != nil || !matched {
		t.Fatalf("%s = %#v, want 40 lowercase hex Git object ID", name, object[name])
	}
}

func validatePinnedHelmObject(t *testing.T, registry map[string]map[string]any, schemaID string, object map[string]any) {
	t.Helper()
	schema, ok := registry[schemaID]
	if !ok {
		t.Fatalf("Helm schema %s is not pinned", schemaID)
	}
	if err := validateSchemaValueWithRegistry(schemaID, schema, object, registry); err != nil {
		t.Fatalf("fixture does not conform to %s: %v", schemaID, err)
	}
}

func assertDigestMatchesBinding(t *testing.T, object map[string]any, bindingField, digestField string) {
	t.Helper()
	binding := objectValue(t, object, bindingField)
	raw, err := canonicalizeJCS(binding)
	if err != nil {
		t.Fatalf("RFC 8785 canonicalize %s: %v", bindingField, err)
	}
	want := fmt.Sprintf("sha256:%x", sha256.Sum256(raw))
	assertStringValue(t, object, digestField, want)
}

func collectSchemaRefs(value any, refs map[string]bool) {
	switch value := value.(type) {
	case map[string]any:
		for key, child := range value {
			if key == "$ref" {
				if ref, ok := child.(string); ok {
					refs[ref] = true
				}
			}
			collectSchemaRefs(child, refs)
		}
	case []any:
		for _, child := range value {
			collectSchemaRefs(child, refs)
		}
	}
}

func assertJCSVector(t *testing.T, name string, value any, want string) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		got, err := canonicalizeJCS(value)
		if err != nil {
			t.Fatalf("canonicalize: %v", err)
		}
		if string(got) != want {
			t.Fatalf("JCS bytes differ\n got: %s\nwant: %s", got, want)
		}
	})
}

func canonicalizeJCS(value any) ([]byte, error) {
	var out strings.Builder
	if err := appendJCS(&out, value); err != nil {
		return nil, err
	}
	return []byte(out.String()), nil
}

func appendJCS(out *strings.Builder, value any) error {
	switch value := value.(type) {
	case nil:
		out.WriteString("null")
	case bool:
		if value {
			out.WriteString("true")
		} else {
			out.WriteString("false")
		}
	case string:
		return appendJCSString(out, value)
	case json.Number:
		number, err := strconv.ParseFloat(string(value), 64)
		if err != nil {
			return fmt.Errorf("unsupported JSON number %q: %w", value, err)
		}
		return appendJCSNumber(out, number)
	case float64:
		return appendJCSNumber(out, value)
	case float32:
		return appendJCSNumber(out, float64(value))
	case int:
		return appendJCSSafeInteger(out, int64(value))
	case int64:
		return appendJCSSafeInteger(out, value)
	case int32:
		return appendJCSSafeInteger(out, int64(value))
	case uint:
		return appendJCSSafeUnsigned(out, uint64(value))
	case uint64:
		return appendJCSSafeUnsigned(out, value)
	case uint32:
		return appendJCSSafeUnsigned(out, uint64(value))
	case []any:
		out.WriteByte('[')
		for index, item := range value {
			if index > 0 {
				out.WriteByte(',')
			}
			if err := appendJCS(out, item); err != nil {
				return err
			}
		}
		out.WriteByte(']')
	case map[string]any:
		keys := make([]string, 0, len(value))
		for key := range value {
			if !utf8.ValidString(key) {
				return fmt.Errorf("object key is not valid UTF-8")
			}
			keys = append(keys, key)
		}
		sort.Slice(keys, func(i, j int) bool { return utf16Less(keys[i], keys[j]) })
		out.WriteByte('{')
		for index, key := range keys {
			if index > 0 {
				out.WriteByte(',')
			}
			if err := appendJCSString(out, key); err != nil {
				return err
			}
			out.WriteByte(':')
			if err := appendJCS(out, value[key]); err != nil {
				return fmt.Errorf("key %q: %w", key, err)
			}
		}
		out.WriteByte('}')
	default:
		return fmt.Errorf("unsupported JCS value type %T", value)
	}
	return nil
}

func appendJCSString(out *strings.Builder, value string) error {
	if !utf8.ValidString(value) {
		return fmt.Errorf("string is not valid UTF-8")
	}
	out.WriteByte('"')
	for _, char := range value {
		switch char {
		case '\b':
			out.WriteString("\\b")
		case '\t':
			out.WriteString("\\t")
		case '\n':
			out.WriteString("\\n")
		case '\f':
			out.WriteString("\\f")
		case '\r':
			out.WriteString("\\r")
		case '"':
			out.WriteString("\\\"")
		case '\\':
			out.WriteString("\\\\")
		default:
			if char >= 0 && char <= 0x1f {
				fmt.Fprintf(out, "\\u%04x", char)
			} else {
				out.WriteRune(char)
			}
		}
	}
	out.WriteByte('"')
	return nil
}

func appendJCSNumber(out *strings.Builder, value float64) error {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return fmt.Errorf("non-finite number is not valid JCS")
	}
	if value == 0 {
		out.WriteByte('0')
		return nil
	}
	abs := math.Abs(value)
	if abs >= 1e-6 && abs < 1e21 {
		out.WriteString(strconv.FormatFloat(value, 'f', -1, 64))
		return nil
	}
	scientific := strconv.FormatFloat(value, 'e', -1, 64)
	parts := strings.SplitN(scientific, "e", 2)
	exponent, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("normalize exponent %q: %w", scientific, err)
	}
	out.WriteString(parts[0])
	out.WriteByte('e')
	if exponent >= 0 {
		out.WriteByte('+')
	}
	out.WriteString(strconv.Itoa(exponent))
	return nil
}

func appendJCSSafeInteger(out *strings.Builder, value int64) error {
	const maxSafeInteger = int64(1<<53 - 1)
	if value < -maxSafeInteger || value > maxSafeInteger {
		return fmt.Errorf("integer %d exceeds the interoperable IEEE-754 range", value)
	}
	out.WriteString(strconv.FormatInt(value, 10))
	return nil
}

func appendJCSSafeUnsigned(out *strings.Builder, value uint64) error {
	const maxSafeInteger = uint64(1<<53 - 1)
	if value > maxSafeInteger {
		return fmt.Errorf("integer %d exceeds the interoperable IEEE-754 range", value)
	}
	out.WriteString(strconv.FormatUint(value, 10))
	return nil
}

func utf16Less(left, right string) bool {
	leftUnits := utf16.Encode([]rune(left))
	rightUnits := utf16.Encode([]rune(right))
	for index := 0; index < len(leftUnits) && index < len(rightUnits); index++ {
		if leftUnits[index] != rightUnits[index] {
			return leftUnits[index] < rightUnits[index]
		}
	}
	return len(leftUnits) < len(rightUnits)
}

func readJSONObject(t *testing.T, path string) map[string]any {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	var object map[string]any
	if err := json.Unmarshal(raw, &object); err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}
	return object
}

func objectValue(t *testing.T, object map[string]any, name string) map[string]any {
	t.Helper()
	value, ok := object[name].(map[string]any)
	if !ok {
		t.Fatalf("%s is %#v, want object", name, object[name])
	}
	return value
}

func assertStringValue(t *testing.T, object map[string]any, name, want string) {
	t.Helper()
	got, ok := object[name].(string)
	if !ok {
		t.Fatalf("%s is %#v, want string", name, object[name])
	}
	if want != "" && got != want {
		t.Errorf("%s = %q, want %q", name, got, want)
	}
}

func assertEqualStringFields(t *testing.T, left, right map[string]any, name string) {
	t.Helper()
	leftValue, leftOK := left[name].(string)
	rightValue, rightOK := right[name].(string)
	if !leftOK || !rightOK || leftValue != rightValue {
		t.Errorf("%s linkage differs: left=%#v right=%#v", name, left[name], right[name])
	}
}

func assertNumberValue(t *testing.T, object map[string]any, name string, want float64) {
	t.Helper()
	got, ok := object[name].(float64)
	if !ok || got != want {
		t.Errorf("%s = %#v, want %v", name, object[name], want)
	}
}

func assertBoolValue(t *testing.T, object map[string]any, name string, want bool) {
	t.Helper()
	got, ok := object[name].(bool)
	if !ok || got != want {
		t.Errorf("%s = %#v, want %v", name, object[name], want)
	}
}

func arrayValue(t *testing.T, object map[string]any, name string) []any {
	t.Helper()
	got, ok := object[name].([]any)
	if !ok {
		t.Fatalf("%s is %#v, want array", name, object[name])
	}
	return got
}

func stringSet(t *testing.T, value any) map[string]bool {
	t.Helper()
	set, ok := stringSetValueChecked(value)
	if !ok {
		t.Fatalf("value is %#v, want string array", value)
	}
	return set
}

func stringSetValue(value any) map[string]bool {
	set, _ := stringSetValueChecked(value)
	return set
}

func stringSetValueChecked(value any) (map[string]bool, bool) {
	items, ok := value.([]any)
	if !ok {
		return nil, false
	}
	set := make(map[string]bool, len(items))
	for _, item := range items {
		text, ok := item.(string)
		if !ok || strings.TrimSpace(text) == "" {
			return nil, false
		}
		set[text] = true
	}
	return set, true
}
