const filterBtn = document.getElementById("filterBtn")
const resetBtn = document.getElementById("resetBtn")
const grid = document.getElementById("artistGrid")

filterBtn.addEventListener("click", function () {
  const minYear = document.getElementById("minYear").value
  const maxYear = document.getElementById("maxYear").value
  const members = document.getElementById("members").value

  const params = new URLSearchParams()
  if (minYear) params.append("minYear", minYear)
  if (maxYear) params.append("maxYear", maxYear)
  if (members) params.append("members", members)

  fetch("/api/filter?" + params.toString())
    .then(function (res) { return res.json() })
    .then(function (artists) {
      grid.innerHTML = ""

      if (!artists || artists.length === 0) {
        grid.innerHTML = "<p class='empty-state'>No artists match your filters.</p>"
        return
      }

      artists.forEach(function (a) {
        const card = document.createElement("a")
        card.href = "/artist/" + a.id
        card.className = "artist-card"
        card.innerHTML = `
          <img src="${a.image}" alt="${a.name}">
          <div class="card-body">
            <div class="card-name">${a.name}</div>
            <div class="card-meta">Est. ${a.creationDate} &middot; ${a.members.length} members</div>
          </div>
        `
        grid.appendChild(card)
      })
    })
})

resetBtn.addEventListener("click", function () {
  document.getElementById("minYear").value = ""
  document.getElementById("maxYear").value = ""
  document.getElementById("members").value = ""
  window.location.reload()
})

  const input = document.getElementById("search")
    const list = document.getElementById("suggestions")

    input.addEventListener("input", function () {
      const query = input.value.trim()

      if (query === "") {
        list.innerHTML = ""
        return
      }

      fetch("/api/search?q=" + encodeURIComponent(query))
        .then(function (res) { return res.json() })
        .then(function (artists) {
          list.innerHTML = ""
          if (!artists || artists.length === 0) return

          artists.forEach(function (a) {
            const li = document.createElement("li")
            const link = document.createElement("a")
            link.href = "/artist/" + a.id
            link.textContent = a.name
            li.appendChild(link)
            list.appendChild(li)
          })
        })
    })