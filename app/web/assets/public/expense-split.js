// an equal split shows in the share boxes what each person will be asked for, read only,
// a custom split starts from those amounts
// htmx runs this file again when a swap renders its component again
(function () {
  // ParseCents in services-expenses.go, 0 for what the server would refuse
  function parseCents(text) {
    var m = /^(\d{1,10})(?:[.,](\d{1,2}))?$/.exec(text.trim());
    return m ? Number(m[1]) * 100 + Number(((m[2] || "") + "00").slice(0, 2)) : 0;
  }

  // splitEqual in services-expenses.go: the leftover cents go to the payer,
  // or one each to the first people when the payer is not one of them
  // shortcut: repeats the Go rule and the repo has no JS tests, move the sums to the server if the rule gets complex
  function splitEqual(total, people, payer) {
    var left = total % people.length;
    var payerAt = people.indexOf(payer);
    return people.map(function (id, i) {
      var share = Math.floor(total / people.length);
      if (payerAt < 0) return i < left ? share + 1 : share;
      return i === payerAt ? share + left : share;
    });
  }

  function update() {
    var form = document.querySelector("[data-expense-split]");
    if (!form) return;
    var equal = form.querySelector('[name="split"][value="equal"]');
    var shares = form.querySelectorAll('[name^="share-"]');
    if (!equal || !equal.checked) {
      shares.forEach(function (input) {
        input.readOnly = false;
      });
      return;
    }
    var people = Array.from(form.querySelectorAll('[name="member"]:checked')).map(function (box) {
      return box.value;
    });
    var total = parseCents(form.querySelector('[name="amount"]').value);
    var amounts = total && people.length ? splitEqual(total, people, form.dataset.expenseSplit) : [];
    var money = new Intl.NumberFormat(document.documentElement.lang, {
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
      useGrouping: false,
    });
    shares.forEach(function (input) {
      var at = people.indexOf(input.name.slice("share-".length));
      input.readOnly = true;
      // empty for nothing owed, a custom split refuses a typed zero and starts from these
      input.value = at >= 0 && amounts[at] ? money.format(amounts[at] / 100) : "";
    });
  }

  if (!window.expenseSplitLoaded) {
    window.expenseSplitLoaded = true;
    document.addEventListener("input", update);
    document.addEventListener("change", update);
  }
  update();
})();
