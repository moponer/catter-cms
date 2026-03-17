// UI工具函数

// 显示加载状态
function showLoading(message = '正在处理，请稍候...') {
	const overlay = document.createElement('div');
	overlay.className = 'loading-overlay';
	overlay.innerHTML = `
		<div class="loading-container">
			<div class="loading-spinner"></div>
			<p class="loading-message">${message}</p>
		</div>
	`;
	document.body.appendChild(overlay);
}

// 隐藏加载状态
function hideLoading() {
	const overlay = document.querySelector('.loading-overlay');
	if (overlay) {
		overlay.remove();
	}
}

// 显示成功提示
function showSuccess(message) {
	const toast = document.createElement('div');
	toast.className = 'toast toast-success';
	toast.innerHTML = `
		<div class="toast-icon">✓</div>
		<div class="toast-message">${message}</div>
	`;
	document.body.appendChild(toast);
	
	setTimeout(() => {
		toast.classList.add('toast-show');
	}, 10);
	
	setTimeout(() => {
		toast.classList.remove('toast-show');
		setTimeout(() => toast.remove(), 300);
	}, 3000);
}

// 显示错误提示
function showError(message) {
	const toast = document.createElement('div');
	toast.className = 'toast toast-error';
	toast.innerHTML = `
		<div class="toast-icon">✗</div>
		<div class="toast-message">${message}</div>
	`;
	document.body.appendChild(toast);
	
	setTimeout(() => {
		toast.classList.add('toast-show');
	}, 10);
	
	setTimeout(() => {
		toast.classList.remove('toast-show');
		setTimeout(() => toast.remove(), 300);
	}, 5000);
}

// 显示警告提示
function showWarning(message) {
	const toast = document.createElement('div');
	toast.className = 'toast toast-warning';
	toast.innerHTML = `
		<div class="toast-icon">⚠</div>
		<div class="toast-message">${message}</div>
	`;
	document.body.appendChild(toast);
	
	setTimeout(() => {
		toast.classList.add('toast-show');
	}, 10);
	
	setTimeout(() => {
		toast.classList.remove('toast-show');
		setTimeout(() => toast.remove(), 300);
	}, 4000);
}

// 显示信息提示
function showInfo(message) {
	const toast = document.createElement('div');
	toast.className = 'toast toast-info';
	toast.innerHTML = `
		<div class="toast-icon">ℹ</div>
		<div class="toast-message">${message}</div>
	`;
	document.body.appendChild(toast);
	
	setTimeout(() => {
		toast.classList.add('toast-show');
	}, 10);
	
	setTimeout(() => {
		toast.classList.remove('toast-show');
		setTimeout(() => toast.remove(), 300);
	}, 3000);
}

// 异步请求包装器
async function fetchWithLoading(url, options = {}) {
	const { loadingMessage = '正在加载...', method = 'GET' } = options;
	
	showLoading(loadingMessage);
	
	try {
		const response = await fetch(url, {
			method,
			headers: {
				'Content-Type': 'application/json',
				'Accept': 'application/json',
			},
			...options,
		});
		
		hideLoading();
		
		if (!response.ok) {
			throw new Error(`HTTP error! status: ${response.status}`);
		}
		
		return await response.json();
	} catch (error) {
		hideLoading();
		showError(`请求失败: ${error.message}`);
		throw error;
	}
}

// 表单提交包装器
async function submitFormWithLoading(form, options = {}) {
	const { loadingMessage = '正在提交...', successMessage = '提交成功' } = options;
	
	showLoading(loadingMessage);
	
	try {
		const formData = new FormData(form);
		const response = await fetch(form.action, {
			method: 'POST',
			body: formData,
			...options,
		});
		
		hideLoading();
		
		if (!response.ok) {
			throw new Error(`HTTP error! status: ${response.status}`);
		}
		
		showSuccess(successMessage);
		return await response.json();
	} catch (error) {
		hideLoading();
		showError(`提交失败: ${error.message}`);
		throw error;
	}
}

// 确认对话框
function confirmAction(message) {
	return new Promise((resolve) => {
		const result = confirm(message);
		resolve(result);
	});
}

// 快捷键支持
document.addEventListener('keydown', function(e) {
	// Ctrl+S 保存
	if (e.ctrlKey && e.key === 's') {
		e.preventDefault();
		const saveButton = document.getElementById('save-button');
		if (saveButton) {
			saveButton.click();
		}
	}
	
	// Ctrl+N 新建
	if (e.ctrlKey && e.key === 'n') {
		e.preventDefault();
		const newButton = document.getElementById('new-button');
		if (newButton) {
			newButton.click();
		}
	}
	
	// ESC 取消
	if (e.key === 'Escape') {
		e.preventDefault();
		const cancelButton = document.getElementById('cancel-button');
		if (cancelButton) {
			cancelButton.click();
		}
	}
	
	// Ctrl+Shift+A 全选
	if (e.ctrlKey && e.shiftKey && e.key === 'A') {
		e.preventDefault();
		const checkboxes = document.querySelectorAll('input[type="checkbox"]');
		checkboxes.forEach(checkbox => checkbox.checked = true);
	}
	
	// Ctrl+Shift+D 删除选中
	if (e.ctrlKey && e.shiftKey && e.key === 'D') {
		e.preventDefault();
		const deleteButton = document.getElementById('batch-delete-button');
		if (deleteButton) {
			deleteButton.click();
		}
	}
});

// 批量操作
function batchDelete(ids) {
	confirmAction(`确定要删除选中的 ${ids.length} 项吗？`).then(confirmed => {
		if (confirmed) {
			fetch('/api/posts/batch', {
				method: 'DELETE',
				headers: {
					'Content-Type': 'application/json',
				},
				body: JSON.stringify({ ids: ids }),
			})
			.then(response => {
				if (response.ok) {
					showSuccess('批量删除成功');
					// 刷新列表
					setTimeout(() => location.reload(), 1000);
				} else {
					showError('批量删除失败');
				}
			})
			.catch(error => {
				showError(`批量删除失败: ${error.message}`);
			});
		}
	});
}

// 批量操作
function batchUpdate(ids, data) {
	confirmAction(`确定要更新选中的 ${ids.length} 项吗？`).then(confirmed => {
		if (confirmed) {
			fetch('/api/posts/batch', {
				method: 'PUT',
				headers: {
					'Content-Type': 'application/json',
				},
				body: JSON.stringify({ ids: ids, data: data }),
			})
			.then(response => {
				if (response.ok) {
					showSuccess('批量更新成功');
					// 刷新列表
					setTimeout(() => location.reload(), 1000);
				} else {
					showError('批量更新失败');
				}
			})
			.catch(error => {
				showError(`批量更新失败: ${error.message}`);
			});
		}
	});
}

// 页面加载完成
window.addEventListener('load', function() {
	// 隐藏加载指示器
	hideLoading();
	
	// 处理Flash消息
	const flashMessages = document.querySelectorAll('.flash-message');
	flashMessages.forEach(message => {
		const messageType = message.dataset.type;
		switch (messageType) {
			case 'success':
				showSuccess(message.textContent);
				break;
			case 'error':
				showError(message.textContent);
				break;
			case 'warning':
				showWarning(message.textContent);
				break;
			case 'info':
				showInfo(message.textContent);
				break;
		}
		message.remove();
	});
});
