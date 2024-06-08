<template>
  <div id="app">
    <AppNavbar />
    <router-view/>
  </div>
</template>

<script>
import AppNavbar from './components/AppNavbar.vue';

export default {
  name: 'App',
  components: {
    AppNavbar,
  },
};
</script>

<style>
body {
  padding-top: 56px; /* Высота navbar */
  margin: 0;
}

#app {
  padding: 0;
  margin: 0;
}
</style>
