(() => {
  const list = document.getElementById("cities");
  const swap = document.getElementById("swap");
  const from = document.getElementById("from");
  const to = document.getElementById("to");
  if (swap && from && to) {
    swap.addEventListener("click", () => {
      [from.value, to.value] = [to.value, from.value];
    });
  }
  if (!list) return;
  let timer;
  const fill = async (input) => {
    const q = input.value.trim();
    if (q.length < 1) return;
    try {
      const res = await fetch("/api/cities?q=" + encodeURIComponent(q));
      if (!res.ok) return;
      const names = await res.json();
      list.innerHTML = names.map((n) => `<option value="${n.replace(/"/g, "&quot;")}">`).join("");
    } catch (_) {}
  };
  [from, to].forEach((input) => {
    if (!input) return;
    input.addEventListener("input", () => {
      clearTimeout(timer);
      timer = setTimeout(() => fill(input), 120);
    });
  });
})();
