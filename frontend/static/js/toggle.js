document.addEventListener('DOMContentLoaded', function() {
    const leftSidebar = document.querySelector('.sidebar-left');
    const rightSidebar = document.querySelector('.sidebar-right');
    const toggleLeft = document.getElementById('toggleLeft');
    const toggleRight = document.getElementById('toggleRight');
    const overlay = document.getElementById('overlay');
    const closeButtons = document.querySelectorAll('.close-sidebar');
    function toggleSidebar(sidebar) {
        const isActive = sidebar.classList.contains('active');
        
        // Close any open sidebar first
        [leftSidebar, rightSidebar].forEach(sb => {
            sb.classList.remove('active');
        });
        overlay.classList.remove('active');
        
        // If we weren't trying to close the sidebar, open it
        if (!isActive) {
            sidebar.classList.add('active');
            overlay.classList.add('active');
        }
    }
    toggleLeft.addEventListener('click', () => toggleSidebar(leftSidebar));
    toggleRight.addEventListener('click', () => toggleSidebar(rightSidebar));
    closeButtons.forEach(button => {
        button.addEventListener('click', () => {
            const sidebarType = button.dataset.sidebar;
            const sidebar = sidebarType === 'left' ? leftSidebar : rightSidebar;
            toggleSidebar(sidebar);
        });
    });
    overlay.addEventListener('click', () => {
        [leftSidebar, rightSidebar].forEach(sidebar => {
            sidebar.classList.remove('active');
        });
        overlay.classList.remove('active');
    });
    // Close sidebars when screen is resized to desktop view
    window.addEventListener('resize', () => {
        if (window.innerWidth >= 1024) {
            [leftSidebar, rightSidebar].forEach(sidebar => {
                sidebar.classList.remove('active');
            });
            overlay.classList.remove('active');
        }
    });
});