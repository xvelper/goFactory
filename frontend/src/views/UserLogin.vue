<template>
  <div class="container mt-5 d-flex justify-content-center">
    <div class="card p-4" style="max-width: 400px; width: 100%;">
      <h2 class="text-center mb-4">Sign in to GitFactory</h2>
      <form @submit.prevent="login">
        <div class="mb-3">
          <label for="username" class="form-label">Username</label>
          <input type="text" v-model="username" class="form-control" id="username" required>
        </div>
        <div class="mb-3">
          <label for="password" class="form-label">Password</label>
          <input type="password" v-model="password" class="form-control" id="password" required>
        </div>
        <button type="submit" class="btn btn-success w-100">Sign in</button>
      </form>
      <hr class="my-4">
      <div class="text-center">
        <p>New to GitFactory? <router-link to="/register" class="text-decoration-none">Create an account</router-link></p>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'UserLogin',
  data() {
    return {
      username: '',
      password: ''
    };
  },
  methods: {
    async login() {
      try {
        const response = await axios.post('https://api.xvelper.ru/api/v1/login', {
          username: this.username,
          password: this.password
        }, { withCredentials: true });

        localStorage.setItem('token', response.data.token);
        const userDetails = await axios.get('https://api.xvelper.ru/api/v1/user_details_jwt', { withCredentials: true });
        localStorage.setItem('username', userDetails.data.username);

        window.dispatchEvent(new Event('storage')); // Dispatch event to notify about the change

        this.$router.push('/my-repositories');
      } catch (error) {
        alert("Введен неправильный логин или пароль!")
        console.error('Error logging in:', error);
      }
    }
  }
};
</script>


<style scoped>
body {
  background-color: #f8f9fa; /* Light background */
  color: #212529; /* Default text color */
}

.form-control {
  background-color: #ffffff; /* White background for form fields */
  border: 1px solid #ced4da; /* Light border */
  color: #212529; /* Default text color */
}

.form-control:focus {
  border-color: #007bff; /* Blue border on focus */
  box-shadow: 0 0 0 0.2rem rgba(0, 123, 255, 0.25); /* Blue shadow on focus */
}

.card {
  background-color: #ffffff; /* White background for the card */
  border: 1px solid #ced4da; /* Light border */
}

.btn-success {
  background-color: #28a745; /* Green button */
  border-color: #28a745; /* Green border */
}

.btn-success:hover {
  background-color: #218838; /* Darker green on hover */
  border-color: #1e7e34; /* Darker green border on hover */
}

.text-muted {
  color: #6c757d !important; /* Muted text color */
}

.text-decoration-none {
  color: #007bff !important; /* Blue link color */
}

.text-decoration-none:hover {
  text-decoration: underline; /* Underline on hover */
}
</style>
