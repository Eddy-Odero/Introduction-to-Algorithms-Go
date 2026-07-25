// Client-side validation + inline error display for forms that don't have
// a backend yet (auth lands in Phase 4, posts in Phase 5). Submitting a
// valid form shows a toast instead of actually posting anywhere.

(function () {
  function showError(field, message) {
    field.classList.add('has-error');
    const msg = field.querySelector('.field-error');
    if (msg && message) msg.textContent = message;
  }

  function clearError(field) {
    field.classList.remove('has-error');
  }

  function validateField(field) {
    const input = field.querySelector('input, textarea');
    if (!input) return true;

    clearError(field);

    if (input.hasAttribute('required') && !input.value.trim()) {
      showError(field, 'This field is required.');
      return false;
    }

    if (input.type === 'email' && input.value && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(input.value)) {
      showError(field, 'Enter a valid email address.');
      return false;
    }

    if (input.type === 'password' && input.value && input.value.length < 8) {
      showError(field, 'Password must be at least 8 characters.');
      return false;
    }

    return true;
  }

  function validateCategoryPicker(form) {
    const picker = form.querySelector('[data-require-one]');
    if (!picker) return true;

    const errorMsg = picker.parentElement.querySelector('.field-error');
    const checked = picker.querySelectorAll('input[type="checkbox"]:checked');

    if (checked.length === 0) {
      picker.classList.add('has-error');
      if (errorMsg) errorMsg.style.display = 'block';
      return false;
    }
    picker.classList.remove('has-error');
    if (errorMsg) errorMsg.style.display = 'none';
    return true;
  }

  function attach(form) {
    form.addEventListener('submit', (e) => {
      e.preventDefault();

      let valid = true;
      form.querySelectorAll('.field').forEach((field) => {
        if (!validateField(field)) valid = false;
      });
      if (!validateCategoryPicker(form)) valid = false;

      if (!valid) {
        window.forumToast('Please fix the highlighted fields.');
        return;
      }

      window.forumToast(form.dataset.successMessage || 'Looks good — this connects to the server in a later phase.');
    });

    // Clear an error as soon as the person starts fixing it.
    form.querySelectorAll('input, textarea').forEach((input) => {
      input.addEventListener('input', () => {
        const field = input.closest('.field');
        if (field) clearError(field);
      });
    });
  }

  document.querySelectorAll('form[data-validate]').forEach(attach);
})();
