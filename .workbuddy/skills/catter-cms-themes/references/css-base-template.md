# Catter CMS base.css 标准模板

> 这是所有主题必须包含的 `base.css` 基础结构。每个主题可以在此基础上扩展，但不得删减 CSS 变量定义部分。

---

## 完整 base.css 模板

```css
/* ============================================
   1. CSS Variables（设计 Token 系统）
   默认值在此定义，模板中通过 <style> 注入用户配置覆盖
   ============================================ */
:root {
    /* ===== 主色系 ===== */
    --primary: #2563eb;
    --primary-hover: #1d4ed8;
    --primary-dark: #0f3460;
    --primary-darker: #16213e;
    --accent: #60a5fa;

    /* ===== 语义色 ===== */
    --success: #10b981;
    --warning: #f59e0b;
    --danger: #ef4444;

    /* ===== 背景色 ===== */
    --bg-primary: #f8fafc;
    --bg-secondary: #ffffff;
    --bg-tertiary: #f1f5f9;

    /* ===== 文字色 ===== */
    --text-primary: #1e293b;
    --text-secondary: #64748b;
    --text-muted: #94a3b8;

    /* ===== 边框 ===== */
    --border-color: #e5e7eb;
    --border-hover: #d1d5db;

    /* ===== 字体 ===== */
    --font-sans: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto,
                 'Helvetica Neue', Arial, 'Noto Sans SC', sans-serif;
    --font-mono: 'SF Mono', Monaco, 'Fira Code', 'Fira Mono', monospace;

    /* ===== 间距 ===== */
    --spacing-xs: 0.25rem;
    --spacing-sm: 0.5rem;
    --spacing-md: 1rem;
    --spacing-lg: 1.5rem;
    --spacing-xl: 2rem;
    --spacing-2xl: 3rem;
    --spacing-3xl: 4rem;

    /* ===== 圆角 ===== */
    --radius-sm: 0.25rem;
    --radius-md: 0.5rem;
    --radius-lg: 0.75rem;
    --radius-xl: 1rem;
    --radius-full: 9999px;

    /* ===== 阴影 ===== */
    --shadow-sm: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
    --shadow-md: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
    --shadow-lg: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);

    /* ===== 布局 ===== */
    --max-width: 1200px;
    --content-width: 760px;

    /* ===== 动画 ===== */
    --transition: all 0.2s ease;
    --transition-slow: all 0.4s ease;
}

/* ============================================
   2. Reset & Base
   ============================================ */
*, *::before, *::after {
    box-sizing: border-box;
    margin: 0;
    padding: 0;
}

html {
    font-size: 16px;
    scroll-behavior: smooth;
    -webkit-text-size-adjust: 100%;
}

body {
    font-family: var(--font-sans);
    font-size: 1rem;
    line-height: 1.7;
    color: var(--text-primary);
    background: var(--bg-primary);
    min-height: 100vh;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
}

a {
    color: var(--accent);
    text-decoration: none;
    transition: color 0.2s;
}
a:hover { color: var(--primary-hover); }

img {
    max-width: 100%;
    height: auto;
    display: block;
}

ul, ol { list-style: none; }

/* ============================================
   3. Header / Navbar
   ============================================ */
.header {
    background: var(--bg-secondary);
    border-bottom: 1px solid var(--border-color);
    position: sticky;
    top: 0;
    z-index: 100;
    box-shadow: var(--shadow-sm);
}

.header-container {
    max-width: var(--max-width);
    margin: 0 auto;
    padding: 0 var(--spacing-xl);
    height: 64px;
    display: flex;
    align-items: center;
    justify-content: space-between;
}

.logo {
    display: flex;
    align-items: center;
    gap: var(--spacing-sm);
    font-size: 1.25rem;
    font-weight: 700;
    color: var(--text-primary);
}
.logo:hover { color: var(--primary); }

.logo-img {
    height: 32px;
    width: auto;
}

.nav-menu {
    display: flex;
    align-items: center;
    gap: var(--spacing-xl);
}

.nav-link {
    font-size: 0.9375rem;
    font-weight: 500;
    color: var(--text-secondary);
    padding: var(--spacing-xs) 0;
    border-bottom: 2px solid transparent;
    transition: var(--transition);
}
.nav-link:hover, .nav-link.active {
    color: var(--primary);
    border-bottom-color: var(--primary);
}

/* ============================================
   4. Footer
   ============================================ */
.footer {
    background: var(--primary);
    color: rgba(255, 255, 255, 0.9);
    padding: var(--spacing-xl) var(--spacing-xl);
    margin-top: var(--spacing-3xl);
}

.footer-container {
    max-width: var(--max-width);
    margin: 0 auto;
    text-align: center;
}

.footer-nav {
    display: flex;
    justify-content: center;
    flex-wrap: wrap;
    gap: var(--spacing-lg);
    margin-bottom: var(--spacing-md);
}

.footer-link {
    color: rgba(255, 255, 255, 0.8);
    font-size: 0.9rem;
    transition: var(--transition);
}
.footer-link:hover { color: white; }

.footer-copy {
    font-size: 0.875rem;
    color: rgba(255, 255, 255, 0.7);
}

/* ============================================
   5. Page Header（页面标题区）
   ============================================ */
.page-header {
    background: linear-gradient(135deg, var(--primary-dark) 0%, var(--primary-darker) 100%);
    color: white;
    padding: var(--spacing-3xl) var(--spacing-xl);
    text-align: center;
}

.page-header h1 {
    font-size: 2.25rem;
    font-weight: 700;
    margin-bottom: var(--spacing-sm);
}

.page-header p {
    font-size: 1.1rem;
    opacity: 0.85;
    max-width: 600px;
    margin: 0 auto;
}

/* ============================================
   6. Empty State（无内容提示）
   ============================================ */
.empty-state {
    text-align: center;
    padding: var(--spacing-3xl) var(--spacing-xl);
    color: var(--text-secondary);
}

.empty-state i {
    font-size: 3rem;
    margin-bottom: var(--spacing-md);
    color: var(--text-muted);
    display: block;
}

/* ============================================
   7. 通用卡片组件
   ============================================ */
.card {
    background: var(--bg-secondary);
    border: 1px solid var(--border-color);
    border-radius: var(--radius-lg);
    box-shadow: var(--shadow-sm);
    transition: box-shadow 0.2s, transform 0.2s;
}
.card:hover {
    box-shadow: var(--shadow-md);
    transform: translateY(-2px);
}

/* ============================================
   8. 按钮
   ============================================ */
.btn-primary {
    display: inline-flex;
    align-items: center;
    gap: var(--spacing-xs);
    padding: var(--spacing-sm) var(--spacing-xl);
    background: var(--primary);
    color: white;
    border-radius: var(--radius-md);
    font-weight: 500;
    transition: var(--transition);
}
.btn-primary:hover {
    background: var(--primary-hover);
    color: white;
    transform: translateY(-1px);
}

/* ============================================
   9. 分页组件
   ============================================ */
.pagination {
    display: flex;
    justify-content: center;
    align-items: center;
    gap: var(--spacing-sm);
    margin-top: var(--spacing-2xl);
    padding-bottom: var(--spacing-2xl);
}

.page-btn {
    min-width: 40px;
    height: 40px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    padding: 0 var(--spacing-sm);
    border: 1px solid var(--border-color);
    border-radius: var(--radius-md);
    color: var(--text-secondary);
    font-size: 0.875rem;
    transition: var(--transition);
}
.page-btn:hover { border-color: var(--primary); color: var(--primary); }
.page-btn.active {
    background: var(--primary);
    border-color: var(--primary);
    color: white;
}

/* ============================================
   10. 侧边栏 Widget
   ============================================ */
.widget {
    background: var(--bg-secondary);
    border: 1px solid var(--border-color);
    border-radius: var(--radius-lg);
    padding: var(--spacing-lg);
    margin-bottom: var(--spacing-lg);
}

.widget-title {
    font-size: 1rem;
    font-weight: 600;
    color: var(--text-primary);
    margin-bottom: var(--spacing-md);
    padding-bottom: var(--spacing-sm);
    border-bottom: 2px solid var(--primary);
}

.widget-category-list li {
    padding: var(--spacing-xs) 0;
    border-bottom: 1px solid var(--border-color);
}
.widget-category-list li:last-child { border-bottom: none; }

.category-link {
    display: flex;
    justify-content: space-between;
    align-items: center;
    color: var(--text-secondary);
    font-size: 0.9rem;
    transition: var(--transition);
}
.category-link:hover { color: var(--primary); padding-left: var(--spacing-sm); }

/* ============================================
   11. 标签云
   ============================================ */
.tags-cloud {
    display: flex;
    flex-wrap: wrap;
    gap: var(--spacing-sm);
}

.tag-item {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    padding: 4px var(--spacing-sm);
    background: var(--bg-tertiary);
    color: var(--text-secondary);
    border: 1px solid var(--border-color);
    border-radius: var(--radius-full);
    font-size: 0.8125rem;
    transition: var(--transition);
}
.tag-item:hover {
    background: var(--primary);
    color: white;
    border-color: var(--primary);
}

.tag-count {
    background: var(--border-color);
    color: var(--text-muted);
    border-radius: var(--radius-full);
    padding: 0 6px;
    font-size: 0.75rem;
}

/* ============================================
   12. 归档列表
   ============================================ */
.archive-year { margin-bottom: var(--spacing-2xl); }
.archive-year-title {
    font-size: 1.5rem;
    font-weight: 700;
    color: var(--primary);
    margin-bottom: var(--spacing-md);
    padding-bottom: var(--spacing-sm);
    border-bottom: 2px solid var(--border-color);
}

.archive-list { padding: 0; }
.archive-item {
    display: flex;
    align-items: baseline;
    gap: var(--spacing-md);
    padding: var(--spacing-sm) 0;
    border-bottom: 1px dashed var(--border-color);
}
.archive-item:last-child { border-bottom: none; }

.archive-date {
    flex-shrink: 0;
    font-size: 0.875rem;
    color: var(--text-muted);
    font-family: var(--font-mono);
    min-width: 50px;
}

.archive-title {
    color: var(--text-primary);
    font-size: 0.9375rem;
    transition: var(--transition);
}
.archive-title:hover { color: var(--primary); }

/* ============================================
   13. 搜索表单
   ============================================ */
.search-form {
    display: flex;
    gap: var(--spacing-sm);
    max-width: 500px;
    margin: var(--spacing-md) auto 0;
}

.search-input {
    flex: 1;
    padding: var(--spacing-sm) var(--spacing-md);
    border: 2px solid rgba(255,255,255,0.4);
    border-radius: var(--radius-full);
    background: rgba(255,255,255,0.15);
    color: white;
    font-size: 0.9375rem;
    outline: none;
    transition: var(--transition);
}
.search-input::placeholder { color: rgba(255,255,255,0.6); }
.search-input:focus { border-color: white; background: rgba(255,255,255,0.25); }

.search-btn {
    padding: var(--spacing-sm) var(--spacing-md);
    background: white;
    color: var(--primary);
    border: none;
    border-radius: var(--radius-full);
    cursor: pointer;
    font-size: 1rem;
    transition: var(--transition);
}
.search-btn:hover { background: var(--bg-primary); transform: scale(1.05); }

.results-info {
    color: var(--text-secondary);
    font-size: 0.9rem;
    margin-bottom: var(--spacing-lg);
    padding: var(--spacing-sm) var(--spacing-md);
    background: var(--bg-secondary);
    border-radius: var(--radius-md);
    border-left: 3px solid var(--primary);
}

/* ============================================
   14. 分类网格
   ============================================ */
.categories-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: var(--spacing-lg);
    padding: var(--spacing-xl) 0;
}

.category-card {
    background: var(--bg-secondary);
    border: 1px solid var(--border-color);
    border-radius: var(--radius-lg);
    padding: var(--spacing-xl);
    text-align: center;
    transition: var(--transition);
}
.category-card:hover {
    border-color: var(--primary);
    box-shadow: var(--shadow-md);
    transform: translateY(-3px);
}

.category-card-name {
    font-size: 1.1rem;
    font-weight: 600;
    color: var(--text-primary);
    margin-bottom: var(--spacing-xs);
}
.category-card-count {
    font-size: 0.875rem;
    color: var(--text-muted);
}

/* ============================================
   15. Read More 链接
   ============================================ */
.read-more {
    display: inline-flex;
    align-items: center;
    gap: var(--spacing-xs);
    color: var(--primary);
    font-size: 0.875rem;
    font-weight: 500;
    margin-top: var(--spacing-sm);
    transition: var(--transition);
}
.read-more:hover { gap: var(--spacing-sm); color: var(--primary-hover); }

/* ============================================
   16. 404 页面
   ============================================ */
.not-found {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: calc(100vh - 64px - 80px);
    padding: var(--spacing-2xl);
    text-align: center;
}

.not-found-code {
    font-size: 8rem;
    font-weight: 900;
    color: var(--primary);
    line-height: 1;
    opacity: 0.3;
}

.not-found-title {
    font-size: 1.5rem;
    color: var(--text-primary);
    margin: var(--spacing-md) 0 var(--spacing-sm);
}

.not-found-desc {
    color: var(--text-secondary);
    margin-bottom: var(--spacing-xl);
}

/* ============================================
   17. Responsive（移动端响应式）
   ============================================ */
@media (max-width: 768px) {
    .header-container {
        flex-direction: column;
        height: auto;
        padding: var(--spacing-md) var(--spacing-md);
        gap: var(--spacing-sm);
    }

    .nav-menu {
        flex-wrap: wrap;
        justify-content: center;
        gap: var(--spacing-md);
    }

    .nav-link {
        font-size: 0.875rem;
    }

    .page-header {
        padding: var(--spacing-2xl) var(--spacing-md);
    }

    .page-header h1 {
        font-size: 1.75rem;
    }

    .categories-grid {
        grid-template-columns: repeat(2, 1fr);
    }

    .archive-item {
        flex-direction: column;
        gap: 2px;
    }

    .pagination {
        flex-wrap: wrap;
    }
}

@media (max-width: 480px) {
    :root { --max-width: 100%; }

    .header-container {
        padding: var(--spacing-sm) var(--spacing-md);
    }

    .page-header h1 {
        font-size: 1.5rem;
    }

    .categories-grid {
        grid-template-columns: 1fr;
    }
}
```
