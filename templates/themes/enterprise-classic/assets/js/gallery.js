/**
 * Enterprise Classic - Gallery JS
 * 画廊页专属：灯箱效果（点击放大、前后翻页、键盘导航）
 */
document.addEventListener('DOMContentLoaded', function() {
    'use strict';

    var gallery = document.querySelector('.gallery-grid');
    if (!gallery) return;

    // 创建灯箱容器
    var lightbox = document.createElement('div');
    lightbox.id = 'gallery-lightbox';
    lightbox.innerHTML = '<div class="lightbox-content"><img src="" alt=""><span class="lightbox-caption"></span><button class="lightbox-close">&times;</button><button class="lightbox-prev"><i class="fa-solid fa-chevron-left"></i></button><button class="lightbox-next"><i class="fa-solid fa-chevron-right"></i></button></div>';
    document.body.appendChild(lightbox);

    var lightboxImg = lightbox.querySelector('img');
    var lightboxCaption = lightbox.querySelector('.lightbox-caption');
    var closeBtn = lightbox.querySelector('.lightbox-close');
    var prevBtn = lightbox.querySelector('.lightbox-prev');
    var nextBtn = lightbox.querySelector('.lightbox-next');

    var items = gallery.querySelectorAll('.gallery-link');
    var currentIndex = 0;

    // 点击图片打开灯箱
    items.forEach(function(item, index) {
        item.addEventListener('click', function(e) {
            e.preventDefault();
            currentIndex = index;
            showImage(index);
            lightbox.classList.add('active');
            document.body.style.overflow = 'hidden';
        });
    });

    function showImage(index) {
        var img = items[index].querySelector('img');
        lightboxImg.src = img.src;
        lightboxImg.alt = img.alt;
        lightboxCaption.textContent = img.alt || '';
    }

    // 关闭灯箱
    closeBtn.addEventListener('click', closeLightbox);
    lightbox.addEventListener('click', function(e) {
        if (e.target === lightbox) closeLightbox();
    });

    function closeLightbox() {
        lightbox.classList.remove('active');
        document.body.style.overflow = '';
    }

    // 上一张/下一张
    prevBtn.addEventListener('click', function() {
        currentIndex = (currentIndex - 1 + items.length) % items.length;
        showImage(currentIndex);
    });

    nextBtn.addEventListener('click', function() {
        currentIndex = (currentIndex + 1) % items.length;
        showImage(currentIndex);
    });

    // 键盘导航
    document.addEventListener('keydown', function(e) {
        if (!lightbox.classList.contains('active')) return;
        if (e.key === 'Escape') closeLightbox();
        if (e.key === 'ArrowLeft') prevBtn.click();
        if (e.key === 'ArrowRight') nextBtn.click();
    });
});
