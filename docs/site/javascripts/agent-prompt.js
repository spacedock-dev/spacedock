// ABOUTME: Copy the agent instruction, including after Material instant navigation.
(() => {
  function enableCopy() {
    if (!navigator.clipboard) return;
    document.querySelectorAll('.sd-prompt-copy').forEach((button) => {
      button.hidden = false;
    });
  }

  document.addEventListener('click', async (event) => {
    const button = event.target.closest('.sd-prompt-copy');
    if (!button || button.disabled) return;
    const prompt = button.closest('.sd-agent-prompt').querySelector('code');
    button.disabled = true;
    try {
      await navigator.clipboard.writeText(prompt.textContent);
      button.textContent = 'copied';
      button.setAttribute('aria-label', 'Agent prompt copied');
      button.classList.add('done');
    } catch {
      button.textContent = 'retry';
      button.setAttribute('aria-label', 'Copy failed. Try again or select the instruction to copy it.');
    }
    setTimeout(() => {
      button.textContent = 'copy';
      button.setAttribute('aria-label', 'Copy agent prompt');
      button.classList.remove('done');
      button.disabled = false;
    }, 1500);
  });

  if (typeof document$ !== 'undefined') document$.subscribe(enableCopy);
  else enableCopy();
})();
