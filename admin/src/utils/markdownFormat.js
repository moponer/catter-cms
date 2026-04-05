/**
 * Markdown 自动格式化工具
 */

/**
 * 格式化 Markdown 内容
 * @param {string} content - 原始 Markdown 内容
 * @returns {string} 格式化后的 Markdown
 */
export function formatMarkdown(content) {
  if (!content) return ''

  let result = content

  // 1. 移除 Windows 换行符
  result = result.replace(/\r\n/g, '\n')

  // 2. 移除连续的多个空行（保留最多2个）
  result = result.replace(/\n{3,}/g, '\n\n')

  // 3. 标题前后添加空行（顶级标题）
  result = result.replace(/\n(#{1,6} )/g, '\n\n$1')

  // 4. 移除标题前面的空行
  result = result.replace(/\n{2,}(#{1,6} )/g, '\n\n$1')

  // 5. 代码块前后添加空行
  result = result.replace(/\n(```)/g, '\n\n$1')
  result = result.replace(/(```)\n/g, '$1\n\n')

  // 6. 列表项前后添加空行（顶级列表）
  result = result.replace(/\n(- |\d+\. |\* )/g, '\n\n$1')

  // 7. 区块引用前后添加空行
  result = result.replace(/\n(> )/g, '\n\n$1')

  // 8. 移除开头和结尾的空行
  result = result.trim()

  // 9. 确保结尾有换行符
  if (!result.endsWith('\n')) {
    result += '\n'
  }

  return result
}

/**
 * 规范化代码块缩进
 * @param {string} content
 * @returns {string}
 */
export function normalizeCodeBlocks(content) {
  // 统一代码块标记
  let result = content

  // 行内代码中的反引号保持不变
  // 代码块统一使用三个反引号
  result = result.replace(/^```(\w*)$/gm, '```')

  return result
}

/**
 * 规范化列表缩进
 * @param {string} content
 * @returns {string}
 */
export function normalizeLists(content) {
  let result = content

  // 无序列表：统一使用 -
  result = result.replace(/^[\*\-] /gm, '- ')

  // 有序列表：确保使用数字和点
  result = result.replace(/^(\d+)\)( )/gm, '$1.$2')

  return result
}

/**
 * 完整格式化 Markdown
 * @param {string} content
 * @returns {string}
 */
export function formatFull(content) {
  if (!content) return ''
  let result = formatMarkdown(content)
  result = normalizeCodeBlocks(result)
  result = normalizeLists(result)
  return result
}

export default formatMarkdown
