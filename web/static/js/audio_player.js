function playAudio(trackId, trackName = '', trackArtist = '', trackImage = '') {
    console.log(`Attempting to play audio for track ID: ${trackId}`);
    
    // Get references to DOM elements
    const persistentPlayer = document.getElementById('persistent-player');
    const audio = document.getElementById('audio');
    const playerTrackName = document.getElementById('player-track-name');
    const playerTrackArtist = document.getElementById('player-track-artist');
    const playerImage = document.getElementById('player-image');
    
    // Create the audio URL
    const audioUrl = `/audio/${trackId}`;
    console.log(`Setting audio src to: ${audioUrl}`);
    
    // Set the audio source
    audio.src = audioUrl;
    
    // Make the player visible and add the body class
    persistentPlayer.style.display = 'flex';
    document.body.classList.add('player-active');
    
    // Update track info if provided
    if (trackName && playerTrackName) playerTrackName.textContent = trackName;
    if (trackArtist && playerTrackArtist) playerTrackArtist.textContent = trackArtist;
    if (trackImage && playerImage) playerImage.src = trackImage;
    
    // Try to find missing info from the closest product card if not provided
    if ((!trackName || !trackArtist || !trackImage) && event && event.target) {
        const productCard = event.target.closest('.product-card');
        if (productCard) {
            const playButton = productCard.querySelector('.play-button');
            
            if (!trackName && playerTrackName) {
                const nameElement = productCard.querySelector('h3');
                if (nameElement) playerTrackName.textContent = nameElement.textContent;
                else if (playButton) playerTrackName.textContent = playButton.getAttribute('data-product-name') || '';
            }
            
            if (!trackArtist && playerTrackArtist) {
                const vendorElement = productCard.querySelector('.vendor a');
                if (vendorElement) playerTrackArtist.textContent = vendorElement.textContent;
                else if (playButton) playerTrackArtist.textContent = playButton.getAttribute('data-product-owner') || '';
            }
            
            if (!trackImage && playerImage) {
                const imageElement = productCard.querySelector('img.product-image');
                if (imageElement) playerImage.src = imageElement.src;
            }
        }
    }
    
    // Save track info to localStorage
    const trackData = {
        id: trackId,
        position: 0,
        isPlaying: true,
        name: playerTrackName ? playerTrackName.textContent : '',
        artist: playerTrackArtist ? playerTrackArtist.textContent : '',
        image: playerImage ? playerImage.src : ''
    };
    localStorage.setItem('currentTrack', JSON.stringify(trackData));
    
    // Set up the play event and handle errors
    audio.oncanplay = () => {
        console.log('Audio can play, attempting to play now');
        let playPromise = audio.play();
        
        if (playPromise !== undefined) {
            playPromise.then(_ => {
                console.log('Audio started playing successfully');
            }).catch(error => {
                console.error('Error playing the audio:', error);
            });
        }
    };
    
    audio.onerror = () => {
        console.error(`Error loading audio from ${audioUrl}`);
        console.error('Audio error code:', audio.error ? audio.error.code : 'unknown');
    };
    
    // Add an ended event to remove the track from localStorage when finished
    audio.onended = () => {
        console.log('Track playback ended');
        // Optional: You can decide whether to close the player or not when track ends
        // For now, just update the localStorage to reflect that it's not playing
        const currentData = JSON.parse(localStorage.getItem('currentTrack') || '{}');
        currentData.isPlaying = false;
        localStorage.setItem('currentTrack', JSON.stringify(currentData));
    };
}

// Persistent audio player functionality
document.addEventListener('DOMContentLoaded', function() {
    // Get DOM references
    const persistentPlayer = document.getElementById('persistent-player');
    const audio = document.getElementById('audio');
    const playerTrackName = document.getElementById('player-track-name');
    const playerTrackArtist = document.getElementById('player-track-artist');
    const playerImage = document.getElementById('player-image');
    const closePlayerBtn = document.getElementById('close-player');
    
    // Add event listener to close the player
    if (closePlayerBtn) {
        closePlayerBtn.addEventListener('click', function() {
            persistentPlayer.style.display = 'none';
            audio.pause();
            document.body.classList.remove('player-active');
            // Clear storage
            localStorage.removeItem('currentTrack');
        });
    }
    
    // Check if there's a saved track in localStorage and restore it
    const savedTrack = localStorage.getItem('currentTrack');
    if (savedTrack) {
        try {
            const trackData = JSON.parse(savedTrack);
            
            // Restore player state
            persistentPlayer.style.display = 'flex';
            document.body.classList.add('player-active');
            
            // Set track info if available
            if (playerTrackName && trackData.name) playerTrackName.textContent = trackData.name;
            if (playerTrackArtist && trackData.artist) playerTrackArtist.textContent = trackData.artist;
            if (playerImage && trackData.image) playerImage.src = trackData.image;
            
            // Set audio source and restore playback position
            if (audio) {
                audio.src = `/audio/${trackData.id}`;
                
                // Set up event when metadata is loaded
                audio.onloadedmetadata = function() {
                    // If there was a saved position, seek to it
                    if (trackData.position) {
                        audio.currentTime = trackData.position;
                    }
                    
                    // If it was playing, resume playback
                    if (trackData.isPlaying) {
                        const playPromise = audio.play();
                        if (playPromise !== undefined) {
                            playPromise.catch(error => {
                                console.error('Error resuming playback:', error);
                            });
                        }
                    }
                };
            }
        } catch (error) {
            console.error('Error restoring player state:', error);
            localStorage.removeItem('currentTrack'); // Clear invalid data
        }
    }
    
    // Save current track data and position every second when playing
    if (audio) {
        setInterval(function() {
            if (!audio.paused && persistentPlayer.style.display !== 'none') {
                const trackData = {
                    id: getCurrentTrackId(),
                    position: audio.currentTime,
                    duration: audio.duration,
                    isPlaying: !audio.paused,
                    name: playerTrackName ? playerTrackName.textContent : '',
                    artist: playerTrackArtist ? playerTrackArtist.textContent : '',
                    image: playerImage ? playerImage.src : ''
                };
                localStorage.setItem('currentTrack', JSON.stringify(trackData));
            }
        }, 1000);
    }
});

// Helper to extract the current track id from the audio src
function getCurrentTrackId() {
    const audio = document.getElementById('audio');
    if (!audio || !audio.src) return null;
    
    // Extract the ID from the audio URL (e.g., /audio/123)
    const matches = audio.src.match(/\/audio\/(\d+)/);
    return matches ? matches[1] : null;
}
