<template>
  <div class="container mt-5">
    <h2>Профиль пользователя</h2>
    <form @submit.prevent="updateProfile">
      <div class="mb-3">
        <label for="firstName" class="form-label">Имя</label>
        <input type="text" class="form-control" id="firstName" v-model="profile.firstName" :disabled="!isEditing">
      </div>
      <div class="mb-3">
        <label for="lastName" class="form-label">Фамилия</label>
        <input type="text" class="form-control" id="lastName" v-model="profile.lastName" :disabled="!isEditing">
      </div>
      <div class="mb-3">
        <label for="email" class="form-label">Почта</label>
        <input type="email" class="form-control" id="email" v-model="profile.email" :disabled="!isEditing">
      </div>
      <div class="mb-3">
        <label for="password" class="form-label">Пароль</label>
        <input type="password" class="form-control" id="password" v-model="profile.password" :disabled="!isEditing">
      </div>
      <button type="button" class="btn btn-primary" v-if="!isEditing" @click="enableEditing">Редактировать</button>
      <button type="submit" class="btn btn-success" v-if="isEditing">Сохранить</button>
      <button type="button" class="btn btn-secondary" v-if="isEditing" @click="cancelEditing">Отмена</button>
    </form>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'UserProfile',
  data() {
    return {
      isEditing: false,
      profile: {
        firstName: '',
        lastName: '',
        email: '',
        password: ''
      }
    };
  },
  methods: {
    enableEditing() {
      this.isEditing = true;
    },
    cancelEditing() {
      this.isEditing = false;
      this.fetchProfile();
    },
    async fetchProfile() {
      try {
        const token = localStorage.getItem('token');
        const response = await axios.get('https://api.xvelper.ru/api/v1/user_profile', {
          headers: {
            'Authorization': `Bearer ${token}`
          },
          withCredentials: true
        });
        this.profile.firstName = response.data.firstName;
        this.profile.lastName = response.data.lastName;
        this.profile.email = response.data.email;
      } catch (error) {
        console.error('Error fetching profile:', error);
      }
    },
    async updateProfile() {
      try {
        const token = localStorage.getItem('token');
        const response = await axios.post('https://api.xvelper.ru/api/v1/update_profile', this.profile, {
          headers: {
            'Authorization': `Bearer ${token}`
          },
          withCredentials: true
        });
        console.log('Profile updated:', response.data);
        this.isEditing = false;
      } catch (error) {
        console.error('Error updating profile:', error);
      }
    }
  },
  created() {
    this.fetchProfile();
  }
};
</script>

<style scoped>
.container {
  max-width: 600px;
}
</style>
