// JavaScript-based infinite scroll pagination
class InfiniteScroll {
  constructor(paginationData) {
    this.currentPage = paginationData.page;
    this.totalPages = paginationData.totalPages;
    this.totalPosts = paginationData.totalPosts;
    this.itemsShown = paginationData.itemsShown;
    this.feedName = paginationData.feedName;
    this.loading = false;
    this.retryCount = 0;
    this.maxRetries = 3;
    
    this.setupObserver();
    
    // Check if we need to load more content initially
    // (when initial content doesn't fill the viewport)
    this.checkInitialFill();
  }
  
  setupObserver() {
    // Create a persistent element at the bottom to observe
    this.createPersistentTrigger();
    
    this.observer = new IntersectionObserver((entries) => {
      if (entries[0].isIntersecting && !this.loading) {
        this.loadMore();
      }
    }, { 
      rootMargin: '100px' 
    });
    
    this.observeTrigger();
  }
  
  createPersistentTrigger() {
    // Create a persistent trigger that won't be replaced
    this.persistentTrigger = document.createElement('div');
    this.persistentTrigger.id = 'js-scroll-trigger';
    this.persistentTrigger.style.height = '1px';
    this.persistentTrigger.style.width = '100%';
    
    // Insert it before the template trigger
    const templateTrigger = document.getElementById('scroll-trigger');
    if (templateTrigger) {
      templateTrigger.parentNode.insertBefore(this.persistentTrigger, templateTrigger);
    } else {
      // Fallback: add to end of feed container
      const feed = document.getElementById('feed');
      if (feed) {
        feed.appendChild(this.persistentTrigger);
      }
    }
  }
  
  observeTrigger() {
    if (this.persistentTrigger && this.observer) {
      this.observer.observe(this.persistentTrigger);
    }
  }
  
  checkInitialFill() {
    // Wait a frame for layout to complete
    requestAnimationFrame(() => {
      this.fillViewportIfNeeded();
    });
  }
  
  fillViewportIfNeeded() {
    // If we're already loading, no more pages, or exceeded retries, stop
    if (this.loading || this.currentPage >= this.totalPages || this.retryCount >= this.maxRetries) return;
    
    // Check if the trigger is visible in the viewport
    // (meaning content doesn't fill the page)
    if (this.isTriggerVisible()) {
      this.loadMore().then((success) => {
        if (success) {
          // Reset retry count on success
          this.retryCount = 0;
          // After loading, check again if we need more
          // Use requestAnimationFrame to wait for DOM update
          requestAnimationFrame(() => {
            this.fillViewportIfNeeded();
          });
        }
      });
    }
  }
  
  isTriggerVisible() {
    if (!this.persistentTrigger) return false;
    
    const rect = this.persistentTrigger.getBoundingClientRect();
    const viewportHeight = window.innerHeight || document.documentElement.clientHeight;
    
    // Check if the trigger is within the viewport (with some margin)
    return rect.top < viewportHeight + 100;
  }
  
  async loadMore() {
    if (this.currentPage >= this.totalPages) return false;
    
    this.loading = true;
    this.showLoading();
    
    try {
      const nextPage = this.currentPage + 1;
      const response = await fetch(`/${this.feedName}/${nextPage}/`);
      
      // Check if response is ok (status 200-299)
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      
      const html = await response.text();
      
      const parser = new DOMParser();
      const doc = parser.parseFromString(html, 'text/html');
      const newItems = doc.querySelectorAll('#feed li');
      const container = document.getElementById('feed');
      
      if (!container) {
        throw new Error('Feed container not found');
      }
      
      newItems.forEach(item => container.appendChild(item));
      
      this.currentPage = nextPage;
      this.itemsShown += newItems.length;
      
      // Update pagination info
      this.updatePaginationInfo();
      
      // Remove our persistent trigger if this was the last page
      if (this.currentPage >= this.totalPages) {
        if (this.persistentTrigger) {
          this.persistentTrigger.remove();
        }
      }
      
      return true;
            
    } catch (error) {
      console.error('Failed to load more content:', error);
      this.retryCount++;
      
      // Show error message if we've exceeded retries
      if (this.retryCount >= this.maxRetries) {
        this.showError('Failed to load more content. Please refresh the page.');
      }
      
      return false;
    } finally {
      this.loading = false;
      this.hideLoading();
    }
  }
  
  showLoading() {
    const indicator = document.querySelector('.loading-indicator');
    if (indicator) indicator.style.display = 'flex';
  }
  
  hideLoading() {
    const indicator = document.querySelector('.loading-indicator');
    if (indicator) indicator.style.display = 'none';
  }
  
  showError(message) {
    const container = document.getElementById('feed');
    if (container) {
      const errorDiv = document.createElement('div');
      errorDiv.className = 'error-message';
      errorDiv.textContent = message;
      container.appendChild(errorDiv);
    }
  }
  
  updatePaginationInfo() {
    const currentPageEl = document.getElementById('current-page');
    const itemsShownEl = document.getElementById('items-shown');
    
    if (currentPageEl) currentPageEl.textContent = this.currentPage;
    if (itemsShownEl) itemsShownEl.textContent = this.itemsShown;
  }
  
  // Clean up observer on page unload
  destroy() {
    if (this.observer) {
      this.observer.disconnect();
    }
  }
}

// Feature detection and initialization
if ('IntersectionObserver' in window && window.paginationData) {
  let infiniteScroll;
  
  // Initialize when DOM is ready
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => {
      infiniteScroll = new InfiniteScroll(window.paginationData);
    });
  } else {
    infiniteScroll = new InfiniteScroll(window.paginationData);
  }
  
  // Clean up on page unload to prevent memory leaks
  window.addEventListener('beforeunload', () => {
    if (infiniteScroll) {
      infiniteScroll.destroy();
    }
  });
}
