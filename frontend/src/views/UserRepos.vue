<template>
    <div class="container mt-5">
      <h2>My Repositories</h2>
      <ul class="list-group">
        <li v-for="repo in repos" :key="repo.id" class="list-group-item">
          {{ repo.name }} ({{ repo.is_public ? 'Public' : 'Private' }})
        </li>
      </ul>
    </div>
  </template>
  
  <script>
  import axios from 'axios';
  
  export default {
    data() {
      return {
        repos: [],
      };
    },
    async created() {
      try {
        const token = localStorage.getItem('token');
        const response = await axios.get('https://api.xvelper.ru/api/v1/user_repos', {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        this.repos = response.data;
      } catch (error) {
        alert('Failed to fetch repositories');
      }
    },
  };
  </script>