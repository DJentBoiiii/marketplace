// Audio Player - Ultra compact version
const AP = (() => {
  // Private vars
  const s = {t: [], i: -1, src: "", v: 1}; // state: tracks, index, source type, last volume
  const e = {}; // Elements cache
  
  // Initialize
  function init() {
    // Cache DOM elements
    ['persistent-player', 'audio', 'player-track-name', 'player-track-artist', 'player-image', 
     'play-pause-btn', 'prev-track-btn', 'next-track-btn', 'close-player', 'progress',
     'current-time', 'duration', 'volume-slider', 'volume-icon']
      .forEach(id => e[id.replace(/-([a-z])/g, (_, c) => c.toUpperCase())] = document.getElementById(id));
    e.pb = document.querySelector('.progress-bar');
    
    // Bind events - one-liners for most events
    if (e.closePlayer) e.closePlayer.onclick = () => (e.persistentPlayer.style.display = 'none', e.audio.pause(),
                                                     document.body.classList.remove('player-active'), localStorage.removeItem('currentTrack'));
    if (e.playPauseBtn) e.playPauseBtn.onclick = () => e.audio.paused ? 
      e.audio.play().then(() => e.playPauseBtn.innerHTML = '<i class="fas fa-pause"></i>') : 
      (e.audio.pause(), e.playPauseBtn.innerHTML = '<i class="fas fa-play"></i>', save());
    if (e.prevTrackBtn) e.prevTrackBtn.onclick = prev;
    if (e.nextTrackBtn) e.nextTrackBtn.onclick = next;
    if (e.pb) e.pb.onclick = e => (e.audio.currentTime = (e.clientX - e.pb.getBoundingClientRect().left) / e.pb.offsetWidth * e.audio.duration, updateUI());
    if (e.volumeSlider) e.volumeSlider.oninput = () => (e.audio.volume = e.volumeSlider.value, updateVolIcon());
    if (e.volumeIcon) e.volumeIcon.onclick = () => e.audio.volume > 0 ?
      (s.v = e.audio.volume, e.audio.volume = e.volumeSlider.value = 0) :
      (e.audio.volume = e.volumeSlider.value = s.v, updateVolIcon());
    
    // Audio events
    if (e.audio) {
      ['timeupdate', 'loadedmetadata'].forEach(ev => e.audio.addEventListener(ev, updateUI));
      e.audio.onended = next;
      e.audio.volume = e.volumeSlider?.value || 1;
      setInterval(save, 5000);
    }
    
    restore();
  }

  // Format time and update UI elements
  function updateUI() {
    if (!e.audio || !e.currentTime || !e.duration || !e.progress) return;
    const fmt = sec => `${~~(sec/60)}:${(~~(sec%60)).toString().padStart(2,'0')}`;
    e.currentTime.textContent = fmt(e.audio.currentTime);
    e.duration.textContent = fmt(e.audio.duration);
    e.progress.style.width = `${(e.audio.currentTime/e.audio.duration)*100}%`;
  }

  // Update volume icon based on level
  function updateVolIcon() {
    if (e.volumeIcon) e.volumeIcon.className = 'fas fa-volume-' + 
      (e.audio.volume === 0 ? 'mute' : e.audio.volume < 0.5 ? 'down' : 'up');
  }

  // Play previous track
  function prev() {
    if (s.t.length <= 1 || s.i <= 0) 
      return e.audio && (e.audio.currentTime = 0, e.audio.play());
    play(...Object.values(s.t[--s.i]));
  }

  // Play next track
  function next() {
    if (s.t.length <= 1 || s.i >= s.t.length - 1) 
      return e.audio && (e.audio.pause(), e.playPauseBtn.innerHTML = '<i class="fas fa-play"></i>', save(false));
    play(...Object.values(s.t[++s.i]));
  }

  // Find track info from DOM context
  function findTracks(id) {
    s.t = []; s.i = -1;
    
    // Check if we're on a playlist or vendor page
    const items = document.querySelectorAll('.playlist-item');
    const cards = document.querySelectorAll('.product-card');
    
    if (items.length > 0) {
      s.src = "playlist";
      collectTracks(items, id);
    } else if (cards.length > 0 && location.href.includes('/vendor/')) {
      s.src = "vendor";
      collectTracks(cards, id);
    } else {
      // Default: single track
      s.src = "single";
      s.t = [{id, name: e.playerTrackName?.textContent || '', 
              owner: e.playerTrackArtist?.textContent || '', 
              image: e.playerImage?.src || ''}];
      s.i = 0;
    }
  }

  // Collect track info from DOM elements
  function collectTracks(items, currentId) {
    items.forEach((item, i) => {
      const btn = item.querySelector('.play-button');
      if (!btn) return;
      
      s.t.push({
        id: btn.getAttribute('data-product-id'),
        name: btn.getAttribute('data-product-name'),
        owner: btn.getAttribute('data-product-owner'),
        image: item.querySelector('.product-image')?.src || ''
      });
      
      if (s.t[s.t.length-1].id === currentId) s.i = s.t.length-1;
    });
  }

  // Save player state to localStorage
  function save(playing = true) {
    if (!e.audio || e.persistentPlayer.style.display === 'none') return;
    
    localStorage.setItem('currentTrack', JSON.stringify({
      id: e.audio.src.match(/\/audio\/(\d+)/)?.[1],
      pos: e.audio.currentTime,
      playing: playing && !e.audio.paused,
      name: e.playerTrackName?.textContent || '',
      artist: e.playerTrackArtist?.textContent || '',
      image: e.playerImage?.src || '',
      src: s.src,
      tracks: s.t,
      idx: s.i
    }));
  }

  // Restore player from localStorage
  function restore() {
    try {
      const d = JSON.parse(localStorage.getItem('currentTrack'));
      if (!d) return;
      
      // Restore state and display player
      if (d.tracks) s.t = d.tracks;
      if (d.idx !== undefined) s.i = d.idx;
      if (d.src) s.src = d.src;
      
      e.persistentPlayer.style.display = 'flex';
      document.body.classList.add('player-active');
      
      // Set track info
      if (e.playerTrackName && d.name) e.playerTrackName.textContent = d.name;
      if (e.playerTrackArtist && d.artist) e.playerTrackArtist.textContent = d.artist;
      if (e.playerImage && d.image) e.playerImage.src = d.image;
      
      // Set audio source and resume
      if (e.audio && d.id) {
        e.audio.src = `/audio/${d.id}`;
        e.audio.onloadedmetadata = () => {
          if (d.pos) e.audio.currentTime = d.pos;
          updateUI();
          
          if (d.playing) {
            e.audio.play().catch(() => {});
            if (e.playPauseBtn) e.playPauseBtn.innerHTML = '<i class="fas fa-pause"></i>';
          } else if (e.playPauseBtn) {
            e.playPauseBtn.innerHTML = '<i class="fas fa-play"></i>';
          }
        };
      }
    } catch (e) {
      localStorage.removeItem('currentTrack');
    }
  }

  // Play a track
  function play(id, name = '', artist = '', image = '') {
    if (!e.audio || !e.persistentPlayer) return;
    
    // Set audio source
    e.audio.src = `/audio/${id}`;
    e.persistentPlayer.style.display = 'flex';
    document.body.classList.add('player-active');
    
    // Update track info if provided
    if (name && e.playerTrackName) e.playerTrackName.textContent = name;
    if (artist && e.playerTrackArtist) e.playerTrackArtist.textContent = artist;
    if (image && e.playerImage) e.playerImage.src = image;
    
    // Try to find missing info from product card
    if ((!name || !artist || !image) && event?.target) {
      const card = event.target.closest('.product-card') || event.target.closest('.playlist-item');
      if (card) {
        const btn = card.querySelector('.play-button');
        if (!name && e.playerTrackName) 
          e.playerTrackName.textContent = card.querySelector('h3')?.textContent || btn?.dataset.productName || '';
        if (!artist && e.playerTrackArtist) 
          e.playerTrackArtist.textContent = card.querySelector('.vendor a')?.textContent || btn?.dataset.productOwner || '';
        if (!image && e.playerImage) 
          e.playerImage.src = card.querySelector('img.product-image')?.src || '';
      }
    }
    
    // Set play icon and find track context
    if (e.playPauseBtn) e.playPauseBtn.innerHTML = '<i class="fas fa-pause"></i>';
    findTracks(id);
    
    // Play audio and save state
    e.audio.oncanplay = () => e.audio.play().then(updateUI).catch(() => {});
    save();
  }

  return {init, play};
})();

// Public interface and auto-init
function playAudio(id, name = '', artist = '', image = '') {
  AP.play(id, name, artist, image);
}

document.addEventListener('DOMContentLoaded', AP.init);
