/**
 * Enterprise Classic - Home JS
 * 首页专属：Hero 轮播、关于轮播、统计动画、新闻 Tab、新闻轮播、滚动渐入、懒加载
 */
(function() {
    'use strict';

    /* ===== Hero 背景轮播 ===== */
    var heroCarousel = document.getElementById('heroBgCarousel');
    if (heroCarousel) {
        var heroSlides = heroCarousel.querySelectorAll('.hero-bg-slide');
        var heroDots = document.querySelectorAll('.hero-carousel-dot');
        if (heroSlides.length > 1) {
            var heroCurrent = 0;
            var heroTimer = null;
            function heroGoTo(idx) {
                heroSlides[heroCurrent].classList.remove('active');
                if (heroDots[heroCurrent]) heroDots[heroCurrent].classList.remove('active');
                heroCurrent = (idx + heroSlides.length) % heroSlides.length;
                heroSlides[heroCurrent].classList.add('active');
                if (heroDots[heroCurrent]) heroDots[heroCurrent].classList.add('active');
            }
            function heroStart() { heroTimer = setInterval(function() { heroGoTo(heroCurrent + 1); }, 5000); }
            function heroStop() { clearInterval(heroTimer); }
            heroDots.forEach(function(dot, i) {
                dot.addEventListener('click', function() { heroStop(); heroGoTo(i); heroStart(); });
            });
            heroStart();
        }
    }

    /* ===== 新闻分类 Tab 切换 ===== */
    var tabBar = document.getElementById('newsTabBar');
    if (tabBar) {
        var tabs = tabBar.querySelectorAll('.news-tab-btn');
        var items = document.querySelectorAll('.news-list-item');
        tabs.forEach(function(tab) {
            tab.addEventListener('click', function() {
                tabs.forEach(function(t) { t.classList.remove('active'); });
                this.classList.add('active');
                var cat = this.getAttribute('data-category');
                items.forEach(function(item) {
                    if (cat === 'all') {
                        item.classList.remove('hidden');
                    } else {
                        var cats = item.getAttribute('data-category') || '';
                        item.classList.toggle('hidden', !cats.split(' ').includes(cat));
                    }
                });
            });
        });
    }

    /* ===== 关于我们图片轮播 ===== */
    var carousel = document.getElementById('aboutCarousel');
    if (carousel) {
        var slides = carousel.querySelectorAll('.about-carousel-slide');
        var dotsContainer = document.getElementById('aboutCarouselDots');
        var prevBtn = document.getElementById('aboutCarouselPrev');
        var nextBtn = document.getElementById('aboutCarouselNext');

        if (slides.length <= 1) {
            if (dotsContainer) dotsContainer.style.display = 'none';
            if (prevBtn) prevBtn.style.display = 'none';
            if (nextBtn) nextBtn.style.display = 'none';
        } else {
            var dots = carousel.querySelectorAll('.about-carousel-dot');
            var current = 0;
            var timer = null;

            function goTo(idx) {
                slides[current].classList.remove('active');
                dots[current].classList.remove('active');
                current = (idx + slides.length) % slides.length;
                slides[current].classList.add('active');
                dots[current].classList.add('active');
            }

            function startTimer() { timer = setInterval(function() { goTo(current + 1); }, 4000); }
            function stopTimer() { clearInterval(timer); }

            if (prevBtn) prevBtn.addEventListener('click', function() { stopTimer(); goTo(current - 1); startTimer(); });
            if (nextBtn) nextBtn.addEventListener('click', function() { stopTimer(); goTo(current + 1); startTimer(); });
            dots.forEach(function(dot, i) {
                dot.addEventListener('click', function() { stopTimer(); goTo(i); startTimer(); });
            });
            carousel.addEventListener('mouseenter', stopTimer);
            carousel.addEventListener('mouseleave', startTimer);
            startTimer();
        }
    }

    /* ===== 新闻轮播自动播放 ===== */
    var newsCarousel = document.getElementById('newsCarousel');
    if (newsCarousel) {
        var newsSlides = newsCarousel.querySelectorAll('.news-slide');
        var newsDots = newsCarousel.querySelectorAll('.carousel-dot');
        var newsCurrent = 0;
        var newsTimer = null;

        function newsGoTo(idx) {
            newsSlides[newsCurrent].classList.remove('active');
            if (newsDots[newsCurrent]) newsDots[newsCurrent].classList.remove('active');
            newsCurrent = (idx + newsSlides.length) % newsSlides.length;
            newsSlides[newsCurrent].classList.add('active');
            if (newsDots[newsCurrent]) newsDots[newsCurrent].classList.add('active');
        }

        function newsStart() { newsTimer = setInterval(function() { newsGoTo(newsCurrent + 1); }, 5000); }
        function newsStop() { clearInterval(newsTimer); }

        if (newsSlides.length > 1) {
            newsDots.forEach(function(dot, i) {
                dot.addEventListener('click', function() { newsStop(); newsGoTo(i); newsStart(); });
            });
            newsCarousel.addEventListener('mouseenter', newsStop);
            newsCarousel.addEventListener('mouseleave', newsStart);
            newsStart();
        }
    }

    /* ===== 数字计数动画 ===== */
    var statValues = document.querySelectorAll('.stat-value[data-target]');
    if (statValues.length > 0) {
        var statsAnimated = false;
        function animateStats() {
            if (statsAnimated) return;
            var statsSection = document.querySelector('.stats-section');
            if (!statsSection) return;
            var rect = statsSection.getBoundingClientRect();
            if (rect.top < window.innerHeight * 0.85) {
                statsAnimated = true;
                statValues.forEach(function(el) {
                    var target = parseInt(el.getAttribute('data-target')) || 0;
                    var duration = 1500;
                    var start = 0;
                    var startTime = null;
                    function step(timestamp) {
                        if (!startTime) startTime = timestamp;
                        var progress = Math.min((timestamp - startTime) / duration, 1);
                        var eased = 1 - Math.pow(1 - progress, 3);
                        el.textContent = Math.floor(eased * target);
                        if (progress < 1) {
                            requestAnimationFrame(step);
                        } else {
                            el.textContent = target;
                        }
                    }
                    requestAnimationFrame(step);
                });
            }
        }
        window.addEventListener('scroll', animateStats);
        animateStats();
    }

    /* ===== 滚动渐入动画 ===== */
    var fadeEls = document.querySelectorAll('.about-section, .services-section, .showcase-section, .news-section, .stats-section');
    if (fadeEls.length > 0) {
        fadeEls.forEach(function(el) {
            el.style.opacity = '0';
            el.style.transform = 'translateY(30px)';
            el.style.transition = 'opacity 0.6s ease, transform 0.6s ease';
        });
        function checkFade() {
            fadeEls.forEach(function(el) {
                var rect = el.getBoundingClientRect();
                if (rect.top < window.innerHeight * 0.88) {
                    el.style.opacity = '1';
                    el.style.transform = 'translateY(0)';
                }
            });
            if (Array.from(fadeEls).every(function(el) {
                return el.style.opacity === '1';
            })) {
                window.removeEventListener('scroll', checkFade);
            }
        }
        window.addEventListener('scroll', checkFade);
        checkFade();
    }

    /* ===== 懒加载图片淡入 ===== */
    document.querySelectorAll('img[loading="lazy"]').forEach(function(img) {
        if (img.complete) {
            img.classList.add('loaded');
        } else {
            img.addEventListener('load', function() { this.classList.add('loaded'); });
            img.addEventListener('error', function() { this.classList.add('loaded'); });
        }
    });
})();
