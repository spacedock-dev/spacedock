# Frontmatter contract

YAML frontmatter stores workflow settings and work item state. Your agent maintains these fields.

| Record | Contents | Schema |
| --- | --- | --- |
| Workflow README | Stages, defaults and decision points | [Workflow schema](https://github.com/spacedock-dev/spacedock/blob/main/docs/schema/workflow-readme.mdschema.yml) |
| Work item | Identity, current stage and outcome | [Entity schema](https://github.com/spacedock-dev/spacedock/blob/main/docs/schema/entity.mdschema.yml) |

The schemas define field names, types and required values. [Workflows and entities](../concepts/workflows-and-entities.md) explains the files you see.

<a id="entity"></a>
<a id="workflow-readme"></a>
