<template>
  <nav class="navbar navbar-expand-lg bg-body-tertiary fixed-top">
    <div class="container">
      <router-link class="navbar-brand" to="/">GoFactory</router-link>
      <button
        class="navbar-toggler"
        type="button"
        data-bs-toggle="collapse"
        data-bs-target="#navbarScroll"
        aria-controls="navbarScroll"
        aria-expanded="false"
        aria-label="Toggle navigation"
      >
        <span class="navbar-toggler-icon"></span>
      </button>
      <div class="collapse navbar-collapse" id="navbarScroll">
        <ul class="navbar-nav me-auto my-2 my-lg-0 navbar-nav-scroll">
          <li class="nav-item">
            <router-link class="nav-link" to="/">Главная</router-link>
          </li>
          <li class="nav-item">
            <router-link class="nav-link" to="/public-repos">Репозитории</router-link>
          </li>
        </ul>
        <form class="d-flex" role="search">
          <input
            class="form-control me-2"
            type="search"
            placeholder="Поиск"
            aria-label="Поиск"
          />
        </form>
        <div v-if="isLoggedIn" class="dropdown">
          <button
            class="btn btn-outline-success dropdown-toggle"
            type="button"
            id="userDropdown"
            data-bs-toggle="dropdown"
            aria-expanded="false"
          >
            {{ username }}
          </button>
          <ul class="dropdown-menu dropdown-menu-end" aria-labelledby="userDropdown">
            <li>
              <router-link class="dropdown-item" to="/profile">My Profile</router-link>
            </li>
            <li>
              <router-link class="dropdown-item" to="/my-repositories">My Repositories</router-link>
            </li>
            <li>
              <hr class="dropdown-divider" />
            </li>
            <li>
              <a class="dropdown-item" href="#" @click="signOut">Sign Out</a>
            </li>
          </ul>
        </div>
        <button v-else class="btn btn-outline-success" @click="login">Login</button>
      </div>
    </div>
  </nav>
</template>

<script>

export default {
  name: 'AppNavbar',
  data() {
    return {
      isLoggedIn: false,
      username: '',
    };
  },
  methods: {
    login() {
      this.$router.push('/login');
    },
    signOut() {
      localStorage.removeItem('token');
      localStorage.removeItem('username');
      this.isLoggedIn = false;
      this.username = '';
      this.$router.push('/');
    },
    checkLogin() {
      const token = localStorage.getItem('token');
      const username = localStorage.getItem('username');
      if (token && username) {
        this.username = username;
        this.isLoggedIn = true;
      }
    },
  },
  created() {
    this.checkLogin();
    window.addEventListener('storage', this.checkLogin); // Listen for changes to localStorage
  },
  beforeUnmount() {
    window.removeEventListener('storage', this.checkLogin);
  }
};
</script>

<style scoped>
.navbar {
  width: 100%;
  max-width: 75%;
  margin: 0 auto;
  z-index: 10;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
}
</style>
