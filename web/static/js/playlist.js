function showPlaylistModal(productId) {
    document.getElementById('modalProductId').value = productId;
    document.getElementById('playlistModal').style.display = 'flex';
}

function closePlaylistModal() {
    document.getElementById('playlistModal').style.display = 'none';
}

window.onclick = function(event) {
    const modal = document.getElementById('playlistModal');
    if (event.target == modal) {
        closePlaylistModal();
    }
}