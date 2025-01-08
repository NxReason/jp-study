const API = {
  radicals: {
    url: '/api/radicals',
    async all() {
      try {
        const res = await fetch(this.url);
        const json = await res.json();
        return json.radicals;
      } catch (err) {
        console.error(err);
        return null;
      }
    },
  },
};

const Radicals = {
  async init() {
    this.$container = document.getElementById('radicals-table-body');
    this.radicals = await API.radicals.all();
    this.populateTable();
  },
  populateTable() {
    console.log(this.radicals);

    const rows = this.radicals.map(this.createRow);
    this.$container.append(...rows);
  },
  createRow(radical) {
    const row = document.createElement('tr');

    const glyph = document.createElement('td');
    glyph.textContent = radical.Glyph;

    const meanings = document.createElement('td');
    if (radical.Meanings) {
      meanings.textContent = radical.Meanings.map(rm => rm.Meaning).join(', ');
    }

    const progress = document.createElement('td');

    const controls = document.createElement('td');
    const editIcon = document.createElement('i');
    editIcon.classList.add('icon');
    const removeIcon = document.createElement('i');
    removeIcon.classList.add('icon');
    controls.append(editIcon, removeIcon);

    row.append(glyph, meanings, progress, controls);
    return row;
  },
};

Radicals.init();
