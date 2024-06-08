import { createApp, reactive } from 'vue';
import App from './App.vue';
import router from './router';
import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootstrap/dist/js/bootstrap.bundle.min.js';


const app = createApp(App);

const globalState = reactive({
  user: JSON.parse(localStorage.getItem('user')) || null
});

app.provide('globalState', globalState);

app.use(router).mount('#app');
