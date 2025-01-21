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

    async update(body) {
      const headers = {
        method: 'PUT',
        'Content-Type': 'application/json',
        body: JSON.stringify(body),
      };

      const res = await fetch(this.url, headers);
      const json = await res.json();
      return json.radical;
    },

    async remove(id) {
      const headers = {
        method: 'DELETE',
        'Content-Type': 'application/json',
        body: JSON.stringify({ id }),
      };

      const res = await fetch(this.url, headers);
      if (!res.ok) {
        console.error(res.statusText);
        return null;
      }
      const json = await res.json();
      return id;
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

  update(radical) {
    this.radicals = this.radicals.map(r => (r.ID === radical.ID ? radical : r));
    console.log(this.radicals);
    this.updateRow(radical);
  },

  remove(id) {
    this.radicals = this.radicals.filter(r => r.ID !== id);
    const row = document.getElementById(`radical-${id}`);
    row.remove();
  },

  populateTable() {
    const rows = this.radicals.map(this.createRow);
    this.$container.append(...rows);
  },
  createRow(radical) {
    const row = document.createElement('tr');
    row.id = `radical-${radical.ID}`;

    const glyph = document.createElement('td');
    glyph.classList.add('glyph-cell');
    glyph.textContent = radical.Glyph;

    const meanings = document.createElement('td');
    meanings.classList.add('meanings-cell');
    if (radical.Meanings) {
      meanings.textContent = radical.Meanings.map(rm => rm.Meaning).join(', ');
    }

    const progress = document.createElement('td');

    const controls = document.createElement('td');
    controls.dataset['id'] = radical.ID;
    controls.classList.add('controls-cell');
    const editIcon = document.createElement('i');
    editIcon.classList.add('icon', 'icon-edit');
    const removeIcon = document.createElement('i');
    removeIcon.classList.add('icon', 'icon-remove');
    controls.append(editIcon, removeIcon);

    row.append(glyph, meanings, progress, controls);
    return row;
  },

  updateRow(radical) {
    const row = document.getElementById(`radical-${radical.ID}`);
    const glyph = row.querySelector('.glyph-cell');
    glyph.textContent = radical.Glyph;

    const meanings = row.querySelector('.meanings-cell');
    if (radical.Meanings) {
      meanings.textContent = radical.Meanings.map(rm => rm.Meaning).join(', ');
    } else {
      meanings.textContent = '';
    }
  },

  setEventHandlers() {
    this.newButton.addEventListener('click', () => {
      Modal.open();
    });

    this.$container.addEventListener('click', async e => {
      const classes = e.target.classList;
      const getId = () => parseInt(e.target.parentNode.dataset['id']);
      switch (true) {
        case classes.contains('icon-remove'):
          const rid = getId();
          const res = await API.radicals.remove(rid);
          if (res) {
            Radicals.remove(rid);
          }
          break;
        case classes.contains('icon-edit'):
          const eid = getId();
          const radical = this.radicals.find(r => r.ID === eid);
          Modal.open(radical);
          break;
      }
    });
  },
};

const Modal = {
  modes: { NEW: 'NEW', UPDATE: 'UPDATE' },
  init() {
    this.container = document.getElementById('modal-radical');
    this.form = document.getElementById('form-radical');
    this.glyphInput = document.getElementById('input-radical-glyph');
    this.namesInput = document.getElementById('input-radical-names');
    this.saveBtn = document.getElementById('btn-save-radical');
    this.closeBtn = document.getElementById('btn-close-modal');
    this.setFormHandlers();

    this.mode = this.modes.NEW;
    this.updateId = null;
  },

  setFormHandlers() {
    this.form.addEventListener('submit', async e => {
      e.preventDefault();
      const radicalData = this.collectData();
      switch (this.mode) {
        case this.modes.NEW:
          const newRadical = await API.radicals.save(radicalData);
          Radicals.add(newRadical);
          break;
        case this.modes.UPDATE:
          const fields = this.collectData();
          const meanings = fields.meanings.map(m => ({
            meaning: m,
          }));
          const updateBody = {
            id: this.updateId,
            glyph: fields.glyph,
            meanings,
          };
          const updatedRadical = await API.radicals.update(updateBody);
          Radicals.update(updatedRadical);
          this.close();
          break;
        default:
          console.error('What to do with this radical?');
          break;
      }
    });
    this.closeBtn.addEventListener('click', e => {
      e.stopPropagation();
      e.preventDefault();
      this.close();
    });
  },

  setFields(radical) {
    this.glyphInput.value = radical.Glyph;
    let names = '';
    if (radical.Meanings?.length)
      names = radical.Meanings.map(m => m.Meaning).join(', ');
    this.namesInput.value = names;
  },

  open(radical) {
    this.setFields(radical ?? { Glyph: '', Meanings: [] });
    if (radical) {
      this.mode = this.modes.UPDATE;
      this.updateId = radical.ID;
    } else {
      this.mode = this.modes.NEW;
      this.updateId = null;
    }
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
