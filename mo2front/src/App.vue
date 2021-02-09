<template>
  <v-app>
    <v-app-bar
      id="appBarElm"
      scroll-target="#scrolling-techniques-6"
      color="primary"
      dark
    >
      <div class="d-flex align-center">
        <v-img
          alt="Vuetify Logo"
          class="shrink mr-2"
          contain
          src="https://cdn.vuetifyjs.com/images/logos/vuetify-logo-dark.png"
          transition="scale-transition"
          width="40"
        />
        <div class="text-h4">MO2</div>
        <!-- <v-img
          alt="Vuetify Name"
          class="shrink mt-1 hidden-sm-and-down"
          contain
          min-width="100"
          src="https://cdn.vuetifyjs.com/images/logos/vuetify-name-dark.png"
          width="100"
        /> -->
      </div>

      <v-spacer></v-spacer>
    </v-app-bar>
    <v-navigation-drawer
      style="z-index: 99999"
      right
      fixed
      bottom
      permanent
      :mini-variant.sync="sideNavVisible"
      expand-on-hover
    >
      <v-list-item class="px-2" :style="`height: ${appBarHeight}px`">
        <v-list-item-avatar>
          <v-img src="https://randomuser.me/api/portraits/men/85.jpg"></v-img>
        </v-list-item-avatar>

        <v-list-item-title>John Leider</v-list-item-title>

        <v-btn icon @click.stop="sideNavVisible = !sideNavVisible">
          <v-icon>mdi-chevron-left</v-icon>
        </v-btn>
      </v-list-item>

      <v-divider></v-divider>

      <v-list dense>
        <v-list-item v-for="item in items" :key="item.title" link>
          <v-list-item-icon>
            <v-icon>{{ item.icon }}</v-icon>
          </v-list-item-icon>

          <v-list-item-content>
            <v-list-item-title>{{ item.title }}</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>

    <v-main>
      <v-parallax
        src="https://cdn.vuetifyjs.com/images/backgrounds/vbanner.jpg"
        height="200"
      >
        <v-row align="center" justify="center">
          <v-col class="text-center" cols="12">
            <h1 class="display-1 font-weight-thin mb-4">MO2</h1>
            <h4 class="subheading">Monkey ❤ Mouse</h4>
          </v-col>
        </v-row>
      </v-parallax>
      <v-container>
        <router-view />
        <!-- <home /> -->
        <v-btn @click="showLogin()">Login</v-btn>
        <account-modal :enable.sync="enable" />
        <v-btn @click="$vuetify.theme.dark = !$vuetify.theme.dark"
          >switch theme</v-btn
        >
        <v-btn @click="sideNavVisible = !sideNavVisible">show side bar</v-btn>
      </v-container>
    </v-main>
  </v-app>
</template>

<script lang="ts">
import Vue from "vue";
import HelloWorld from "./components/HelloWorld.vue";
import AccountModal from "./components/AccountModal.vue";
import Vuelidate from "vuelidate";
import Component from "vue-class-component";
import Home from "./views/Home.vue";
// import "bulma/bulma.sass";
Vue.use(Vuelidate);

@Component({
  components: {
    AccountModal,
    HelloWorld,
    Home,
  },
})
export default class App extends Vue {
  enable = false;
  sideNavVisible = true;
  items = [
    { title: "Home", icon: "mdi-home-city" },
    { title: "My Account", icon: "mdi-account" },
    { title: "Users", icon: "mdi-account-group-outline" },
  ];
  created() {
    window.addEventListener("resize", () => {
      setTimeout(() => {
        this.appBarHeight = document.getElementById("appBarElm").clientHeight;
      }, 500);
    });
  }
  mounted() {
    this.$vuetify.application.right = 58;
    this.appBarHeight = document.getElementById("appBarElm").clientHeight;
  }
  showLogin() {
    this.enable = true;
  }
  appBarHeight = 64;
}
</script>
<style>
</style>