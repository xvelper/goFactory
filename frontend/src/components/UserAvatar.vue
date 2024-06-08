<template>
  <div>
    <img
      v-if="isLoggedIn"
      :src="avatarUrl"
      @click="toggleSidebar"
      class="user-avatar"
      alt="User Avatar"
    />
    <button v-else @click="login" class="btn btn-outline-success">Login</button>
  </div>
</template>

<script>
export default {
  data() {
    return {
      isLoggedIn: false,
      avatarUrl: '', // URL аватара пользователя
    };
  },
  methods: {
    login() {
      this.$router.push('/login');
    },
    toggleSidebar() {
      this.$emit('toggle-sidebar');
    },
  },
  created() {
    // Здесь вы можете проверить, авторизован ли пользователь, и установить аватар
    const token = localStorage.getItem('token');
    if (token) {
      this.isLoggedIn = true;
      // Установите URL аватара пользователя
      this.avatarUrl = 'https://via.placeholder.com/40'; // Пример заглушки
    }
  },
};
</script>

<style scoped>
.user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  cursor: pointer;
}
</style>
