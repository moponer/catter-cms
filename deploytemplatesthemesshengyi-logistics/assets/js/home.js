/**
 * shengyi-logistics 首页 JS
 * Hero 背景轮播（自动切换 + 导航点交互）
 */
(function () {
    'use strict';

    var carousel = document.getElementById('heroBgCarousel');
    var dotsContainer = document.getElementById('heroCarouselDots');
    if (!carousel) return;

    var slides = carousel.querySelectorAll('.hero-bg-slide');
    if (slides.length <= 1) return; // 单张图不需要轮播

    var dots = dotsContainer ? dotsContainer.querySelectorAll('.hero-carousel-dot') : [];
    var current = 0;
    var interval = 5000; // 5秒一轮
    var timer = null;

    function goTo(index) {
        if (index === current) return;
        slides[current].classList.remove('active');
        if (dots[current]) dots[current].classList.remove('active');

        current = index;

        slides[current].classList.add('active');
        if (dots[current]) dots[current].classList.add('active');
    }

    function next() {
        goTo((current + 1) % slides.length);
    }

    // 自动播放
    function start() {
        stop();
        timer = setInterval(next, interval);
    }

    function stop() {
        if (timer) { clearInterval(timer); timer = null; }
    }

    // 导航点点击
    for (var i = 0; i < dots.length; i++) {
        (function (idx) {
            dots[idx].addEventListener('click', function () {
                goTo(idx);
                start(); // 重置计时器
            });
        })(i);
    }

    // 鼠标悬停暂停
    carousel.addEventListener('mouseenter', stop);
    carousel.addEventListener('mouseleave', start);

    // 启动
    start();
})();
