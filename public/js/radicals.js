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

    async save(body) {
      const headers = {
        method: 'POST',
        'Content-Type': 'application/json',
        body: JSON.stringify(body),
      };

      const res = await fetch(this.url, headers);
      const json = await res.json();
      return json.radical;
    },
  },
};

const Radicals = {
  async init() {
    this.$container = document.getElementById('radicals-table-body');
    this.radicals = await API.radicals.all();
    this.populateTable();

    this.newButton = document.getElementById('btn-new-radical');
    this.setEventHandlers();
  },

  add(radical) {
    this.radicals.push(radical);
    const row = this.createRow(radical);
    this.$container.append(row);
  },

  populateTable() {
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

  setEventHandlers() {
    this.newButton.addEventListener('click', () => {
      Modal.open();
    });
  },
};

const Modal = {
  init() {
    this.container = document.getElementById('modal-radical');

    this.form = document.getElementById('form-radical');
    this.glyphInput = document.getElementById('input-radical-glyph');
    this.namesInput = document.getElementById('input-radical-names');
    this.saveBtn = document.getElementById('btn-save-radical');
    this.closeBtn = document.getElementById('btn-close-modal');
    this.setFormHandlers();
  },

  setFormHandlers() {
    this.form.addEventListener('submit', async e => {
      e.preventDefault();
      const radicalData = this.collectData();
      const radical = await API.radicals.save(radicalData);
      Radicals.add(radical);
      this.close();
    });
    this.closeBtn.addEventListener('click', e => {
      e.stopPropagation();
      e.preventDefault();
      this.close();
    });
  },

  open() {
    this.container.style.display = 'block';
  },
  close() {
    this.container.style.display = 'none';
  },

  collectData() {
    const glyph = this.glyphInput.value;
    let meanings = this.namesInput.value.split(',').map(m => m.trim());
    return { glyph, meanings };
  },
};

Radicals.init();
Modal.init();
