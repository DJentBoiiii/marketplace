function showPlaylistModal(productId) {
    document.getElementById('modalProductId').value = productId;
    document.getElementById('newPlaylistProductId').value = productId;
    
    // Fetch the playlists
    fetch('/playlist/get-user-playlists')
        .then(response => response.json())
        .then(data => {
            const selectElement = document.querySelector('#playlistModal select[name="playlist_id"]');
            // Clear any existing options except the first one
            while (selectElement.options.length > 1) {
                selectElement.remove(1);
            }
            
            // Add the playlists to the select element
            data.forEach(playlist => {
                const option = document.createElement('option');
                option.value = playlist.Id;
                option.textContent = `${playlist.Name} (${playlist.ItemCount} треків)`;
                selectElement.appendChild(option);
            });
            
            // Display the modal
            document.getElementById('playlistModal').style.display = 'flex';
        })
        .catch(error => {
            console.error('Error fetching playlists:', error);
            // Still show the modal even if there's an error
            document.getElementById('playlistModal').style.display = 'flex';
        });
}

function closePlaylistModal() {
    document.getElementById('playlistModal').style.display = 'none';
}

function switchTab(tabId) {
    // Hide all tab contents
    const tabContents = document.getElementsByClassName('tab-content');
    for (let i = 0; i < tabContents.length; i++) {
        tabContents[i].style.display = 'none';
    }
    
    // Remove active class from all tab buttons
    const tabButtons = document.getElementsByClassName('tab-button');
    for (let i = 0; i < tabButtons.length; i++) {
        tabButtons[i].classList.remove('active');
    }
    
    // Show the selected tab content
    document.getElementById(tabId).style.display = 'block';
    
    // Add active class to the clicked tab button
    const activeTabButton = document.querySelector(`[onclick="switchTab('${tabId}')"]`);
    activeTabButton.classList.add('active');
}

window.onclick = function(event) {
    const modal = document.getElementById('playlistModal');
    if (event.target == modal) {
        closePlaylistModal();
    }
}