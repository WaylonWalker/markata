// JavaScript-based infinite scroll pagination
class InfiniteScroll {
  constructor(paginationData) {
    this.currentPage = paginationData.page;
    this.totalPages = paginationData.totalPages;
    this.totalPosts = paginationData.totalPosts;
    this.itemsShown = paginationData.itemsShown;
    this.feedName = paginationData.feedName;
    this.loading = false;
    
    this.setupObserver();
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
  
  async loadMore() {
    if (this.currentPage >= this.totalPages) return;
    
    this.loading = true;
    this.showLoading();
    
    try {
      const nextPage = this.currentPage + 1;
      const response = await fetch(`/${this.feedName}/${nextPage}/`);
      const html = await response.text();
      
      const parser = new DOMParser();
      const doc = parser.parseFromString(html, 'text/html');
      const newItems = doc.querySelectorAll('#feed li');
      const container = document.getElementById('feed');
      
      newItems.forEach(item => container.appendChild(item));
      
      this.currentPage = nextPage;
      this.itemsShown += newItems.length;
      
      // Update pagination info
      this.updatePaginationInfo();
      
      // Update URL for bookmarkability
      if (history.pushState) {
        history.pushState({}, '', `/${this.feedName}/${nextPage}/`);
      }
      
      // Re-observe trigger for next page
      if (this.currentPage < this.totalPages) {
        // No need to re-observe, our persistent trigger stays in place
        // Just ensure it's still being observed
        this.observeTrigger();
      } else {
        // Remove our persistent trigger if this was the last page
        if (this.persistentTrigger) {
          this.persistentTrigger.remove();
        }
      }
            
    } catch (error) {
      console.error('Failed to load more content:', error);
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
  
  updatePaginationInfo() {
    const currentPageEl = document.getElementById('current-page');
    const itemsShownEl = document.getElementById('items-shown');
    
    if (currentPageEl) currentPageEl.textContent = this.currentPage;
    if (itemsShownEl) itemsShownEl.textContent = this.itemsShown;
  }
}

// Feature detection and initialization
if ('IntersectionObserver' in window && window.paginationData) {
  // Initialize when DOM is ready
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => {
      new InfiniteScroll(window.paginationData);
    });
  } else {
    new InfiniteScroll(window.paginationData);
  }
}