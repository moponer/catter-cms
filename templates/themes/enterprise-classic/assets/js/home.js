/* Enterprise Classic v3 - Home Page JS */
(function() {
    'use strict';

    /* ===== Banner 背景轮播 ===== */
    (function initBannerCarousel() {
        var track = document.getElementById('bannerBgCarousel');
        if (!track) return;
        var slides = track.querySelectorAll('.banner-bg-slide');
        var dots = document.getElementById('bannerDots');
        if (!dots || slides.length <= 1) return;
        var dotEls = dots.querySelectorAll('.banner-dot');
        var current = 0;
        var interval = 5000;
        var timer = null;

        function showSlide(idx) {
            slides[current].classList.remove('active');
            dotEls[current].classList.remove('active');
            current = idx;
            if (current >= slides.length) current = 0;
            if (current < 0) current = slides.length - 1;
            slides[current].classList.add('active');
            dotEls[current].classList.add('active');
        }

        function nextSlide() { showSlide(current + 1); }

        dotEls.forEach(function(dot, i) {
            dot.addEventListener('click', function() { showSlide(i); resetTimer(); });
        });

        function resetTimer() {
            clearInterval(timer);
            timer = setInterval(nextSlide, interval);
        }
        resetTimer();
    })();

    /* ===== About 轮播 ===== */
    (function initAboutCarousel() {
        var track = document.getElementById('aboutCarousel');
        if (!track) return;
        var slides = track.querySelectorAll('.about-carousel-slide');
        if (slides.length <= 1) return;
        var dotsContainer = document.getElementById('aboutCarouselDots');
        var dotEls = dotsContainer ? dotsContainer.querySelectorAll('.about-carousel-dot') : [];
        var prevBtn = document.getElementById('aboutCarouselPrev');
        var nextBtn = document.getElementById('aboutCarouselNext');
        var current = 0;

        function showSlide(idx) {
            slides[current].classList.remove('active');
            if (dotEls[current]) dotEls[current].classList.remove('active');
            current = idx;
            if (current >= slides.length) current = 0;
            if (current < 0) current = slides.length - 1;
            slides[current].classList.add('active');
            if (dotEls[current]) dotEls[current].classList.add('active');
        }

        if (prevBtn) prevBtn.addEventListener('click', function() { showSlide(current - 1); });
        if (nextBtn) nextBtn.addEventListener('click', function() { showSlide(current + 1); });

        dotEls.forEach(function(dot, i) {
            dot.addEventListener('click', function() { showSlide(i); });
        });
    })();

    /* ===== 服务项目卡片轮播（左右箭头） ===== */
    (function initServicesCarousel() {
        var wrapper = document.getElementById('servicesCarousel');
        if (!wrapper) return;
        var grid = wrapper.querySelector('.services-grid');
        if (!grid) return;
        var cards = grid.querySelectorAll('.service-card-item');
        if (cards.length <= 4) {
            /* 不需要轮播，隐藏箭头 */
            var prev = document.getElementById('servicesPrev');
            var next = document.getElementById('servicesNext');
            if (prev) prev.style.display = 'none';
            if (next) next.style.display = 'none';
            return;
        }
        var prevBtn = document.getElementById('servicesPrev');
        var nextBtn = document.getElementById('servicesNext');
        var position = 0;
        var cardWidth = cards[0].offsetWidth + 24; /* gap */
        var visibleCount = 4;
        var maxPosition = cards.length - visibleCount;

        function updatePosition() {
            grid.style.transform = 'translateX(-' + (position * cardWidth) + 'px)';
        }

        if (prevBtn) prevBtn.addEventListener('click', function() {
            if (position > 0) { position--; updatePosition(); }
        });

        if (nextBtn) nextBtn.addEventListener('click', function() {
            if (position < maxPosition) { position++; updatePosition(); }
        });
    })();

    /* ===== 新闻Tab切换 ===== */
    (function initNewsTabs() {
        var tabBtns = document.querySelectorAll('.section-tab-group .tab-btn');
        if (!tabBtns.length) return;

        tabBtns.forEach(function(btn) {
            btn.addEventListener('click', function() {
                tabBtns.forEach(function(b) { b.classList.remove('active'); });
                btn.classList.add('active');
                /* Tab 切换逻辑可扩展：根据 data-tab 显示/隐藏不同内容 */
            });
        });
    })();

    /* ===== 滚动渐入动画 ===== */
    (function initFadeInUp() {
        var elements = document.querySelectorAll('.fade-in-up');
        if (!elements.length || !IntersectionObserver) {
            elements.forEach(function(el) { el.classList.add('visible'); });
            return;
        }
        var observer = new IntersectionObserver(function(entries) {
            entries.forEach(function(entry) {
                if (entry.isIntersecting) {
                    entry.target.classList.add('visible');
                    observer.unobserve(entry.target);
                }
            });
        }, { threshold: 0.12, rootMargin: '0px 0px -40px 0px' });
        elements.forEach(function(el) { observer.observe(el); });
    })();

})();
