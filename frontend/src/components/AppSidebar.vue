<template>
  <div v-if="sidebarOpen" class="sidebar">
    <div class="sidebar-header">
      <h3>{{ username }}</h3>
      <button @click="toggleSidebar">X</button>
    </div>
    <ul class="list-unstyled">
      <li><router-link to="/profile">My Profile</router-link></li>
      <li><router-link to="/my-repositories">My Repositories</router-link></li>
      <li><a href="#" @click="signOut">Sign Out</a></li>
    </ul>
  </div>
</template>

<script>
export default {
  name: 'AppSidebar',
  props: {
    sidebarOpen: {
      type: Boolean,
      required: true
    }
  },
  data() {
    return {
      username: '',
    };
  },
  methods: {
    toggleSidebar() {
      this.$emit('update:sidebarOpen', !this.sidebarOpen);
    },
    signOut() {
      localStorage.removeItem('token');
      this.$router.push('/');
      this.toggleSidebar();
    },
  },
  created() {
    const token = localStorage.getItem('token');
    if (token) {
      this.username = 'User';
    }
  },
};
</script>

<style scoped>
.sidebar {
  position: fixed;
  top: 0;
  right: 0;
  width: 250px;
  height: 100%;
  background: #343a40;
  color: #fff;
  padding: 15px;
  z-index: 1000;
  box-shadow: -2px 0 5px rgba(0, 0, 0, 0.5);
}
.sidebar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.sidebar-header h3 {
  margin: 0;
}
.sidebar ul {
  list-style: none;
  padding: 0;
}
.sidebar ul li {
  margin: 10px 0;
}
.sidebar ul li a {
  color: #fff;
  text-decoration: none;
}
.sidebar ul li a:hover {
  text-decoration: underline;
}
</style>
