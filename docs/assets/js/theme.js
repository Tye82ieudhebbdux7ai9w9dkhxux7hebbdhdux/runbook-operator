(function() {
  'use strict';

  // Theme management
  const ThemeManager = {
    STORAGE_KEY: 'runbook-operator-theme',
    THEMES: {
      LIGHT: 'light',
      DARK: 'dark',
      AUTO: 'auto'
    },

    init() {
      this.loadTheme();
      this.setupToggle();
      this.setupSystemThemeListener();
      this.setupMobileNav();
    },

    loadTheme() {
      const savedTheme = localStorage.getItem(this.STORAGE_KEY);
      const systemPrefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;

      let theme;
      if (savedTheme && savedTheme !== this.THEMES.AUTO) {
        theme = savedTheme;
      } else {
        theme = systemPrefersDark ? this.THEMES.DARK : this.THEMES.LIGHT;
      }

      this.applyTheme(theme);
    },

    applyTheme(theme) {
      document.documentElement.setAttribute('data-theme', theme);
      this.updateToggleButton(theme);

      // Update meta theme-color for mobile browsers
      this.updateMetaThemeColor(theme);
    },

    updateMetaThemeColor(theme) {
      let metaThemeColor = document.querySelector('meta[name="theme-color"]');
      if (!metaThemeColor) {
        metaThemeColor = document.createElement('meta');
        metaThemeColor.name = 'theme-color';
        document.head.appendChild(metaThemeColor);
      }

      metaThemeColor.content = theme === this.THEMES.DARK ? '#161b22' : '#ffffff';
    },

    updateToggleButton(theme) {
      const button = document.getElementById('theme-toggle-btn');
      if (button) {
        button.setAttribute('aria-label',
          theme === this.THEMES.DARK ? 'Switch to light mode' : 'Switch to dark mode'
        );
      }
    },

    toggleTheme() {
      const currentTheme = document.documentElement.getAttribute('data-theme');
      const newTheme = currentTheme === this.THEMES.DARK ? this.THEMES.LIGHT : this.THEMES.DARK;

      this.applyTheme(newTheme);
      localStorage.setItem(this.STORAGE_KEY, newTheme);

      // Add animation class for smooth transition
      document.body.classList.add('theme-transitioning');
      setTimeout(() => {
        document.body.classList.remove('theme-transitioning');
      }, 300);
    },

    setupToggle() {
      const toggleButton = document.getElementById('theme-toggle-btn');
      if (toggleButton) {
        toggleButton.addEventListener('click', () => this.toggleTheme());

        // Add keyboard support
        toggleButton.addEventListener('keydown', (e) => {
          if (e.key === 'Enter' || e.key === ' ') {
            e.preventDefault();
            this.toggleTheme();
          }
        });
      }
    },

    setupSystemThemeListener() {
      const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
      mediaQuery.addEventListener('change', (e) => {
        const savedTheme = localStorage.getItem(this.STORAGE_KEY);
        if (!savedTheme || savedTheme === this.THEMES.AUTO) {
          const newTheme = e.matches ? this.THEMES.DARK : this.THEMES.LIGHT;
          this.applyTheme(newTheme);
        }
      });
    },

    setupMobileNav() {
      const navTrigger = document.querySelector('.nav-trigger');
      const navMenu = document.querySelector('.trigger');

      if (navTrigger && navMenu) {
        navTrigger.addEventListener('click', () => {
          navMenu.classList.toggle('active');
        });

        // Close mobile nav when clicking outside
        document.addEventListener('click', (e) => {
          if (!navTrigger.contains(e.target) && !navMenu.contains(e.target)) {
            navMenu.classList.remove('active');
          }
        });

        // Close mobile nav when pressing escape
        document.addEventListener('keydown', (e) => {
          if (e.key === 'Escape') {
            navMenu.classList.remove('active');
          }
        });
      }
    }
  };

  // Code copy functionality
  const CodeCopyManager = {
    init() {
      this.addCopyButtons();
    },

    addCopyButtons() {
      const codeBlocks = document.querySelectorAll('pre code, .highlight code');

      codeBlocks.forEach((codeBlock) => {
        const pre = codeBlock.closest('pre') || codeBlock.closest('.highlight');
        if (pre && !pre.querySelector('.copy-button')) {
          this.createCopyButton(pre, codeBlock);
        }
      });
    },

    createCopyButton(container, codeBlock) {
      const button = document.createElement('button');
      button.className = 'copy-button';
      button.innerHTML = `
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
          <path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1"></path>
        </svg>
        Copy
      `;
      button.title = 'Copy code to clipboard';

      button.addEventListener('click', () => this.copyCode(button, codeBlock));

      container.style.position = 'relative';
      container.appendChild(button);
    },

    async copyCode(button, codeBlock) {
      const code = codeBlock.textContent || codeBlock.innerText;

      try {
        await navigator.clipboard.writeText(code);
        this.showCopySuccess(button);
      } catch (err) {
        // Fallback for older browsers
        this.fallbackCopy(code, button);
      }
    },

    fallbackCopy(text, button) {
      const textArea = document.createElement('textarea');
      textArea.value = text;
      textArea.style.position = 'fixed';
      textArea.style.left = '-999999px';
      textArea.style.top = '-999999px';
      document.body.appendChild(textArea);
      textArea.focus();
      textArea.select();

      try {
        document.execCommand('copy');
        this.showCopySuccess(button);
      } catch (err) {
        console.error('Failed to copy code:', err);
      }

      document.body.removeChild(textArea);
    },

    showCopySuccess(button) {
      const originalContent = button.innerHTML;
      button.innerHTML = `
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="20,6 9,17 4,12"></polyline>
        </svg>
        Copied!
      `;
      button.classList.add('copied');

      setTimeout(() => {
        button.innerHTML = originalContent;
        button.classList.remove('copied');
      }, 2000);
    }
  };

  // Smooth scrolling for anchor links
  const ScrollManager = {
    init() {
      this.setupSmoothScrolling();
      this.setupScrollToTop();
    },

    setupSmoothScrolling() {
      document.querySelectorAll('a[href^="#"]').forEach(anchor => {
        anchor.addEventListener('click', (e) => {
          const target = document.querySelector(anchor.getAttribute('href'));
          if (target) {
            e.preventDefault();
            target.scrollIntoView({
              behavior: 'smooth',
              block: 'start'
            });
          }
        });
      });
    },

    setupScrollToTop() {
      // Create scroll to top button
      const scrollButton = document.createElement('button');
      scrollButton.className = 'scroll-to-top';
      scrollButton.innerHTML = `
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="12" y1="19" x2="12" y2="5"></line>
          <polyline points="5,12 12,5 19,12"></polyline>
        </svg>
      `;
      scrollButton.title = 'Scroll to top';
      scrollButton.setAttribute('aria-label', 'Scroll to top');

      scrollButton.addEventListener('click', () => {
        window.scrollTo({ top: 0, behavior: 'smooth' });
      });

      document.body.appendChild(scrollButton);

      // Show/hide scroll button based on scroll position
      window.addEventListener('scroll', () => {
        if (window.pageYOffset > 300) {
          scrollButton.classList.add('visible');
        } else {
          scrollButton.classList.remove('visible');
        }
      });
    }
  };

  // Initialize everything when DOM is ready
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => {
      ThemeManager.init();
      CodeCopyManager.init();
      ScrollManager.init();
    });
  } else {
    ThemeManager.init();
    CodeCopyManager.init();
    ScrollManager.init();
  }
})();

