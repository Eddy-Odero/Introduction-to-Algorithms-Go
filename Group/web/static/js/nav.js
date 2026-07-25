// Shared across every page: mobile nav toggle + a tiny toast helper.
// Kept dependency-free — vanilla DOM only.

(function () {
  const burger = document.querySelector('.nav-burger');
  const panel = document.querySelector('.nav-mobile-panel');

  if (burger && panel) {
    burger.addEventListener('click', () => {
      const isOpen = panel.classList.toggle('open');
      burger.setAttribute('aria-expanded', String(isOpen));
    });

    // Close the mobile panel if the viewport grows back past the breakpoint.
    window.addEventListener('resize', () => {
      if (window.innerWidth > 720) {
        panel.classList.remove('open');
        burger.setAttribute('aria-expanded', 'false');
      }
    });
  }
})();

// window.forumToast(message) — shows a small bottom toast, used by
// reactions.js and forms.js so every page gets the same feedback style.
window.forumToast = function forumToast(message) {
  let toast = document.querySelector('.toast');
  if (!toast) {
    toast = document.createElement('div');
    toast.className = 'toast';
    document.body.appendChild(toast);
  }
  toast.textContent = message;
  // Restart the animation even if a toast is already showing.
  toast.classList.remove('show');
  // Force reflow so the class removal/addition isn't batched by the browser.
  void toast.offsetWidth;
  toast.classList.add('show');

  clearTimeout(toast._hideTimer);
  toast._hideTimer = setTimeout(() => toast.classList.remove('show'), 2200);
};
