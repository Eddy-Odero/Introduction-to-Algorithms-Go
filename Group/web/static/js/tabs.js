// Profile page tab switching. Markup contract:
//   <div class="tabs">
//     <button class="tab active" data-panel="my-posts">my posts</button>
//     <button class="tab" data-panel="liked-posts">liked posts</button>
//   </div>
//   <div class="tab-panel" id="my-posts">...</div>
//   <div class="tab-panel" id="liked-posts" hidden>...</div>

(function () {
  const tabs = document.querySelectorAll('.tab[data-panel]');
  if (!tabs.length) return;

  tabs.forEach((tab) => {
    tab.addEventListener('click', () => {
      tabs.forEach((t) => t.classList.remove('active'));
      tab.classList.add('active');

      const targetId = tab.dataset.panel;
      document.querySelectorAll('.tab-panel').forEach((panel) => {
        panel.hidden = panel.id !== targetId;
      });
    });
  });
})();
