// Like/Dislike toggling for posts and comments.
//
// Markup contract (see index.html / post.html for real examples):
//   <div class="action-row" data-target="post-1">
//     <button class="action-btn" data-type="like" data-count="14">...</button>
//     <button class="action-btn" data-type="dislike" data-count="2">...</button>
//   </div>
//
// `data-count` is the baseline count with NO reaction applied. The actual
// like/dislike backend doesn't exist until Phase 7 — this only simulates
// the interaction locally (persisted per-browser via localStorage) so the
// page feels real while we build the rest of the frontend.

(function () {
  const STORAGE_KEY = 'forum_reactions_v1';

  function loadState() {
    try {
      return JSON.parse(localStorage.getItem(STORAGE_KEY)) || {};
    } catch {
      return {};
    }
  }

  function saveState(state) {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(state));
  }

  function applyRowState(row, reaction) {
    const likeBtn = row.querySelector('[data-type="like"]');
    const dislikeBtn = row.querySelector('[data-type="dislike"]');
    if (!likeBtn || !dislikeBtn) return;

    const likeBase = parseInt(likeBtn.dataset.count, 10) || 0;
    const dislikeBase = parseInt(dislikeBtn.dataset.count, 10) || 0;

    likeBtn.classList.toggle('liked', reaction === 'like');
    dislikeBtn.classList.toggle('disliked', reaction === 'dislike');

    const likeCount = likeBase + (reaction === 'like' ? 1 : 0);
    const dislikeCount = dislikeBase + (reaction === 'dislike' ? 1 : 0);

    likeBtn.querySelector('.action-count').textContent = likeCount;
    dislikeBtn.querySelector('.action-count').textContent = dislikeCount;
  }

  function init() {
    const state = loadState();
    const rows = document.querySelectorAll('.action-row[data-target]');

    rows.forEach((row) => {
      const target = row.dataset.target;
      applyRowState(row, state[target]);

      row.querySelectorAll('[data-type="like"], [data-type="dislike"]').forEach((btn) => {
        btn.addEventListener('click', () => {
          const type = btn.dataset.type;
          const current = state[target];
          const next = current === type ? null : type; // clicking again removes it

          if (next) {
            state[target] = next;
          } else {
            delete state[target];
          }
          saveState(state);
          applyRowState(row, state[target]);
        });
      });
    });
  }

  document.addEventListener('DOMContentLoaded', init);
})();
