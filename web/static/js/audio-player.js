function playAudio(trackId) {
    const audioPlayer = document.getElementById('audio-player');
    const audio = document.getElementById('audio');
    
    audio.src = `/audio/${trackId}`;
    audioPlayer.style.display = 'flex';
    audio.oncanplay = () => audio.play();
}
