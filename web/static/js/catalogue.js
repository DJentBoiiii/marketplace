// catalogue.js - Handles interactions for catalogue pages
document.addEventListener('DOMContentLoaded', function() {
    console.log('Catalogue page initialized');
    
    // Initialize playlist modal functionality
    initPlaylistModal();
    
    // Use a global click handler for play buttons
    document.addEventListener('click', function(event) {
        // Check if the clicked element or any of its parents is a play button
        const playButton = event.target.closest('.play-button');
        
        if (playButton) {
            event.preventDefault();
            event.stopPropagation();
            
            console.log('Play button clicked through global handler');
            const productId = playButton.getAttribute('data-product-id');
            
            if (productId) {
                console.log('Playing audio for product ID:', productId);
                // Make sure audio_player.js is loaded and the function is available
                if (typeof playAudio === 'function') {
                    playAudio(productId);
                } else {
                    console.error('playAudio function not found. Is audio_player.js loaded?');
                }
            } else {
                console.error('Missing product ID on play button');
            }
        }
    });
    
    // Log that all play buttons are available
    const buttons = document.querySelectorAll('.play-button');
    console.log(`Found ${buttons.length} play buttons on the page`);
    buttons.forEach((btn, index) => {
        console.log(`Button ${index + 1} ID:`, btn.getAttribute('data-product-id'));
    });
});

// Initialize playlist modal functionality
function initPlaylistModal() {
    window.showPlaylistModal = function(productId) {
        const modal = document.getElementById('playlistModal');
        if (!modal) return;
        
        document.getElementById('modalProductId').value = productId;
        document.getElementById('newPlaylistProductId').value = productId;
        modal.style.display = 'flex';
    };
    
    window.closePlaylistModal = function() {
        const modal = document.getElementById('playlistModal');
        if (modal) {
            modal.style.display = 'none';
        }
    };
    
    window.switchTab = function(tabId) {
        const tabs = document.querySelectorAll('.tab-content');
        const buttons = document.querySelectorAll('.tab-button');
        
        // Hide all tabs
        tabs.forEach(tab => {
            tab.style.display = 'none';
        });
        
        // Deactivate all buttons
        buttons.forEach(button => {
            button.classList.remove('active');
        });
        
        // Show the selected tab
        document.getElementById(tabId).style.display = 'block';
        
        // Activate the button that opened the tab
        const activeButton = document.querySelector(`.tab-button[onclick*="${tabId}"]`);
        if (activeButton) {
            activeButton.classList.add('active');
        }
    };
}

// Update the player information if UI elements exist
function updatePlayerInfo(productId, productName, productOwner, productImage) {
    // Check if there are player UI elements to update
    const playerInfo = document.getElementById('player-info');
    const playerTrackName = document.getElementById('player-track-name');
    const playerTrackArtist = document.getElementById('player-track-artist');
    const playerImage = document.getElementById('player-image');
    
    if (playerInfo) {
        playerInfo.style.display = 'flex';
    }
    
    if (playerTrackName) {
        playerTrackName.textContent = productName;
    }
    
    if (playerTrackArtist) {
        playerTrackArtist.textContent = productOwner;
    }
    
    if (playerImage && productImage) {
        playerImage.src = productImage;
        playerImage.alt = productName;
    }
}