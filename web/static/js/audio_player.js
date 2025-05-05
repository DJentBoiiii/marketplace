function playAudio(trackId) {
    console.log(`Attempting to play audio for track ID: ${trackId}`);
    const audioPlayer = document.getElementById('audio-player');
    const audio = document.getElementById('audio');
    
    // Create the audio URL
    const audioUrl = `/audio/${trackId}`;
    console.log(`Setting audio src to: ${audioUrl}`);
    
    // Set the audio source
    audio.src = audioUrl;
    
    // Make the audio player visible
    audioPlayer.style.display = 'flex';
    
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
        console.error('Audio error code:', audio.error.code);
    };
}
