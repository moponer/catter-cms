/**
 * Enterprise Classic v3 - Main JS
 * 全局功能：移动端菜单、滚动动画、轮播、Tab 切换
 */
(function() {
    'use strict';

    /* ===== 移动端菜单 ===== */
    var toggle = document.getElementById('headerToggle');
    var menu = document.getElementById('headerMenu');
    if (toggle && menu) {
        toggle.addEventListener('click', function() {
            menu.classList.toggle('open');
            toggle.classList.toggle('active');
        });
    }

    /* ===== 滚动渐入动画 ===== */
    var fadeElements = document.querySelectorAll('.fade-in-up');
    if (fadeElements.length > 0 && 'IntersectionObserver' in window) {
        var observer = new IntersectionObserver(function(entries) {
            entries.forEach(function(entry) {
                if (entry.isIntersecting) {
                    entry.target.classList.add('visible');
                }
            });
        }, { threshold: 0.1, rootMargin: '0px 0px -50px 0px' });
        fadeElements.forEach(function(el) { observer.observe(el); });
    } else {
        fadeElements.forEach(function(el) { el.classList.add('visible'); });
    }

    /* ===== 数据统计计数动画 ===== */
    var statValues = document.querySelectorAll('.stat-value');
    if (statValues.length > 0 && 'IntersectionObserver' in window) {
        var statsObserver = new IntersectionObserver(function(entries) {
            entries.forEach(function(entry) {
                if (entry.isIntersecting) {
                    var el = entry.target;
                    var target = parseInt(el.getAttribute('data-target'), 10) || 0;
                    animateNumber(el, target);
                    statsObserver.unobserve(el);
                }
            });
        }, { threshold: 0.5 });
        statValues.forEach(function(el) { statsObserver.observe(el); });
    }

    function animateNumber(el, target) {
        var duration = 1500;
        var start = 0;
        var startTime = null;
        function step(timestamp) {
            if (!startTime) startTime = timestamp;
            var progress = Math.min((timestamp - startTime) / duration, 1);
            var easeProgress = 1 - Math.pow(1 - progress, 3);
            el.textContent = Math.floor(easeProgress * target);
            if (progress < 1) {
                requestAnimationFrame(step);
            } else {
                el.textContent = target;
            }
        }
        requestAnimationFrame(step);
    }

    /* ===== Banner 背景轮播 ===== */
    var bannerCarousel = document.getElementById('bannerBgCarousel');
    if (bannerCarousel) {
        var bannerSlides = bannerCarousel.querySelectorAll('.banner-bg-slide');
        var bannerDots = document.querySelectorAll('.banner-dot');
        if (bannerSlides.length > 1) {
            var bannerIndex = 0;
            setInterval(function() {
                bannerSlides[bannerIndex].classList.remove('active');
                if (bannerDots[bannerIndex]) bannerDots[bannerIndex].classList.remove('active');
                bannerIndex = (bannerIndex + 1) % bannerSlides.length;
                bannerSlides[bannerIndex].classList.add('active');
                if (bannerDots[bannerIndex]) bannerDots[bannerIndex].classList.add('active');
            }, 5000);
            bannerDots.forEach(function(dot, i) {
                dot.addEventListener('click', function() {
                    bannerSlides[bannerIndex].classList.remove('active');
                    if (bannerDots[bannerIndex]) bannerDots[bannerIndex].classList.remove('active');
                    bannerIndex = i;
                    bannerSlides[bannerIndex].classList.add('active');
                    if (bannerDots[bannerIndex]) bannerDots[bannerIndex].classList.add('active');
                });
            });
        }
    }

    /* ===== 关于我们轮播 ===== */
    var aboutCarousel = document.getElementById('aboutCarousel');
    if (aboutCarousel) {
        var aboutSlides = aboutCarousel.querySelectorAll('.about-carousel-slide');
        var aboutDots = document.querySelectorAll('.about-carousel-dot');
        var aboutPrev = document.getElementById('aboutCarouselPrev');
        var aboutNext = document.getElementById('aboutCarouselNext');
        var aboutIndex = 0;

        function gotoAboutSlide(idx) {
            aboutSlides[aboutIndex].classList.remove('active');
            if (aboutDots[aboutIndex]) aboutDots[aboutIndex].classList.remove('active');
            aboutIndex = idx;
            aboutSlides[aboutIndex].classList.add('active');
            if (aboutDots[aboutIndex]) aboutDots[aboutIndex].classList.add('active');
        }

        if (aboutSlides.length > 1) {
            setInterval(function() {
                var next = (aboutIndex + 1) % aboutSlides.length;
                gotoAboutSlide(next);
            }, 6000);
            aboutDots.forEach(function(dot, i) {
                dot.addEventListener('click', function() { gotoAboutSlide(i); });
            });
            if (aboutPrev) {
                aboutPrev.addEventListener('click', function() {
                    var prev = (aboutIndex - 1 + aboutSlides.length) % aboutSlides.length;
                    gotoAboutSlide(prev);
                });
            }
            if (aboutNext) {
                aboutNext.addEventListener('click', function() {
                    var next = (aboutIndex + 1) % aboutSlides.length;
                    gotoAboutSlide(next);
                });
            }
        }
    }

    /* ===== 产品展示 Tab 切换 ===== */
    var showcaseTabs = document.getElementById('showcaseTabs');
    var showcaseGrid = document.getElementById('showcaseGrid');
    if (showcaseTabs && showcaseGrid) {
        var tabBtns = showcaseTabs.querySelectorAll('.showcase-tab-btn');
        var showcaseCards = showcaseGrid.querySelectorAll('.showcase-card');
        tabBtns.forEach(function(btn) {
            btn.addEventListener('click', function() {
                var cat = this.getAttribute('data-category');
                tabBtns.forEach(function(b) { b.classList.remove('active'); });
                this.classList.add('active');
                showcaseCards.forEach(function(card) {
                    var cardCat = card.getAttribute('data-category');
                    card.style.display = (cat === '全部' || cardCat === cat) ? '' : 'none';
                });
            });
        });
    }

    /* ===== 图片懒加载 ===== */
    var lazyImages = document.querySelectorAll('img[loading="lazy"]');
    lazyImages.forEach(function(img) {
        if (img.complete) {
            img.classList.add('loaded');
        } else {
            img.addEventListener('load', function() { this.classList.add('loaded'); });
        }
    });
})();
