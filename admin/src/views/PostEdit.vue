<template>
  <div class="post-edit-container">
    <div class="page-header">
      <h1>{{ isEdit ? '编辑文章' : '新建文章' }}</h1>
      <el-button type="primary" @click="savePost" :loading="loading"
        >保存</el-button
      >
    </div>
    <el-form
      :model="postForm"
      :rules="rules"
      ref="postFormRef"
      label-width="100px"
    >
      <el-form-item label="标题" prop="title">
        <el-input
          v-model="postForm.title"
          placeholder="请输入文章标题"
        ></el-input>
      </el-form-item>
      <el-form-item label="slug" prop="slug">
        <el-input
          v-model="postForm.slug"
          placeholder="请输入文章别名"
        ></el-input>
      </el-form-item>
      <el-form-item label="分类" prop="category_id">
        <el-select v-model="postForm.category_id" placeholder="请选择分类">
          <el-option
            v-for="category in categories"
            :key="category.id"
            :label="category.name"
            :value="category.id"
          ></el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="状态" prop="status">
        <el-select v-model="postForm.status" placeholder="请选择状态">
          <el-option label="草稿" value="draft"></el-option>
          <el-option label="已发布" value="published"></el-option>
          <el-option label="已归档" value="archived"></el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="摘要" prop="excerpt">
        <el-input
          v-model="postForm.excerpt"
          type="textarea"
          placeholder="请输入文章摘要"
        ></el-input>
      </el-form-item>
      <el-form-item label="内容" prop="content">
        <el-input
          v-model="postForm.content"
          type="textarea"
          :rows="10"
          placeholder="请输入文章内容"
        ></el-input>
      </el-form-item>
      <el-form-item label="特色图片">
        <el-upload
          class="upload-demo"
          action="/api/media/upload"
          :on-success="handleImageUpload"
          :before-upload="beforeImageUpload"
        >
          <el-button type="primary">上传图片</el-button>
          <template #tip>
            <div class="el-upload__tip">
              只能上传 JPG/PNG 文件，且不超过 2MB
            </div>
          </template>
        </el-upload>
        <img
          v-if="postForm.featured_image"
          :src="postForm.featured_image"
          class="featured-image"
        />
      </el-form-item>
      <el-form-item label="SEO标题">
        <el-input
          v-model="postForm.seo_title"
          placeholder="请输入SEO标题"
        ></el-input>
      </el-form-item>
      <el-form-item label="SEO描述">
        <el-input
          v-model="postForm.seo_description"
          type="textarea"
          placeholder="请输入SEO描述"
        ></el-input>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { usePostStore, useCategoryStore } from '../store';

const route = useRoute();
const router = useRouter();
const postStore = usePostStore();
const categoryStore = useCategoryStore();

const postFormRef = ref(null);
const loading = ref(false);
const isEdit = ref(false);
const postId = ref('');

// 表单数据
const postForm = reactive({
  title: '',
  slug: '',
  content: '',
  excerpt: '',
  featured_image: '',
  status: 'draft',
  category_id: '',
  seo_title: '',
  seo_description: '',
});

// 验证规则
const rules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  slug: [{ required: true, message: '请输入别名', trigger: 'blur' }],
  content: [{ required: true, message: '请输入内容', trigger: 'blur' }],
};

// 分类列表
const categories = ref([]);

onMounted(async () => {
  // 加载分类
  await categoryStore.fetchCategories();
  categories.value = categoryStore.categories;

  // 检查是否是编辑模式
  postId.value = route.params.id;
  if (postId.value) {
    isEdit.value = true;
    await loadPostData();
  }
});

const loadPostData = async () => {
  await postStore.fetchPost(postId.value);
  const post = postStore.currentPost;
  if (post) {
    Object.assign(postForm, post);
  }
};

const savePost = async () => {
  if (!postFormRef.value) return;

  const valid = await postFormRef.value.validate();
  if (!valid) return;

  loading.value = true;

  let success;
  if (isEdit.value) {
    success = await postStore.updatePost(postId.value, postForm);
  } else {
    success = await postStore.createPost(postForm);
  }

  if (success) {
    ElMessage.success(isEdit.value ? '更新成功' : '创建成功');
    router.push('/posts');
  } else {
    ElMessage.error(isEdit.value ? '更新失败' : '创建失败');
  }

  loading.value = false;
};

const handleImageUpload = (response) => {
  postForm.featured_image = response.data.path;
  ElMessage.success('图片上传成功');
};

const beforeImageUpload = (file) => {
  const isJPG = file.type === 'image/jpeg' || file.type === 'image/png';
  const isLt2M = file.size / 1024 / 1024 < 2;

  if (!isJPG) {
    ElMessage.error('只能上传 JPG/PNG 文件');
  }
  if (!isLt2M) {
    ElMessage.error('文件大小不能超过 2MB');
  }

  return isJPG && isLt2M;
};
</script>

<style scoped>
.post-edit-container {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header h1 {
  color: #303133;
}

.featured-image {
  max-width: 200px;
  margin-top: 10px;
  border-radius: 4px;
}
</style>
