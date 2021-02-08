<template>
  <div class="page-container">
    <account-modal></account-modal>
  </div>
  <!-- <div id="app">
    <div id="nav">
      <router-link to="/">Home</router-link> |
      <router-link to="/about">About</router-link>
    </div>
    <router-view/>
  </div> -->
</template>
<script lang="ts">
import Vue from "vue";
import "vue-material/dist/vue-material.min.css";
import "vue-material/dist/theme/default.css";
import "bulma/bulma.sass";
import VueMaterial from "vue-material";
import Component from "vue-class-component";
import AccountModal from "./components/AccountModal.vue";
import Vuelidate from "vuelidate";
Vue.use(Vuelidate);
Vue.use(VueMaterial);
@Component({
  components: {
    AccountModal,
  },
})
export default class App extends Vue {}
</script>

<style lang="scss">
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}

#nav {
  padding: 30px;

  a {
    font-weight: bold;
    color: #2c3e50;

    &.router-link-exact-active {
      color: #42b983;
    }
  }
}
</style>
