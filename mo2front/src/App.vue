<template>
  <v-app>
    <v-app-bar
      id="appBarElm"
      scroll-target="#scrolling-techniques-6"
      color="primary"
      dark
      app
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
      <v-app-bar-nav-icon
        v-if="$vuetify.breakpoint.smAndDown"
        @click.stop="
          () => {
            drawer = !drawer;
          }
        "
      ></v-app-bar-nav-icon>
    </v-app-bar>
    <v-navigation-drawer
      style="z-index: 99999"
      right
      fixed
      :expand-on-hover="this.$vuetify.breakpoint.mdAndUp"
      v-model="drawerProp"
    >
      <v-list-item class="px-2" :style="`height: ${appBarHeight}px`">
        <avatar :user="user" />

        <v-list-item-title>{{
          isUser ? user.name : "未登录"
        }}</v-list-item-title>

        <v-btn icon v-if="!isUser" @click="showLogin()"> 登录 </v-btn>
      </v-list-item>

      <v-divider></v-divider>

      <v-list dense>
        <v-list-item
          v-for="item in items"
          :key="item.title"
          :to="item.href"
          v-show="item.show"
        >
          <v-list-item-icon>
            <v-icon>{{ item.icon }}</v-icon>
          </v-list-item-icon>

          <v-list-item-content>
            <v-list-item-title>{{ item.title }}</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
      <template v-slot:append>
        <v-list-item>
          <v-list-item-icon>
            <v-icon>mdi-theme-light-dark</v-icon>
          </v-list-item-icon>

          <v-list-item-content>
            <v-list-item-title>
              <v-switch
                dense
                :hide-details="true"
                class="pl-3 ma-0"
                @change="changeTheme"
                v-model="$vuetify.theme.dark"
                label="dark mode"
              ></v-switch>
            </v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item>
          <v-list-item-icon>
            <v-icon>mdi-account-cog</v-icon>
          </v-list-item-icon>

          <v-list-item-content>
            <v-list-item-title>Settings</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </template>
    </v-navigation-drawer>
    <v-main>
      <router-view :user="user" />
      <account-modal :enable.sync="enable" />
      <!-- <v-btn @click="showLogin()">Login</v-btn>
      <account-modal :enable.sync="enable" />
      <v-btn @click="$vuetify.theme.dark = !$vuetify.theme.dark"
        >switch theme</v-btn
      >
      <v-btn @click="sideNavVisible = !sideNavVisible">show side bar</v-btn> -->
    </v-main>
  </v-app>
</template>

<script lang="ts">
import Vue from "vue";
import AccountModal from "./components/AccountModal.vue";
import Vuelidate from "vuelidate";
import Component from "vue-class-component";
import { User } from "./models";
import { GetInitials } from "./utils";
import Avatar from "./components/UserAvatar.vue";
// import "bulma/bulma.sass";
Vue.use(Vuelidate);

@Component({
  components: {
    AccountModal,
    Avatar,
  },
})
export default class App extends Vue {
  drawer = false;
  user: User = {
    name: "leezeeyee",
    email: "easilylazy@mo2.com",
    description: "all work no play made me a dull girl",
    site: "www.mo2.pro",
    createTime: "2020-1-1",
    id: "xxxxxxxxxxxxxxx",
    avatar: "",
    roles: ["user"],
  };
  enable = false;
  items = [
    { title: "Home", icon: "mdi-home-city", href: "/", show: true },
    {
      title: "My Home",
      icon: "mdi-account",
      href: "/account",
      show: this.isUser,
    },
    { title: "About", icon: "mdi-alpha-a-circle", href: "/about", show: true },
  ];
  get isUser() {
    return this.user.roles.length > 0;
  }
  get initials(): string {
    return GetInitials(this.user.name);
  }

  created() {
    try {
      this.$vuetify.theme.dark = JSON.parse(
        localStorage.getItem("darkTheme")
      ) as boolean;
    } catch (err) {}
    window.addEventListener("resize", () => {
      setTimeout(() => {
        this.onResize(false);
      }, 500);
    });
  }

  get drawerProp(): boolean {
    return this.drawer || this.$vuetify.breakpoint.mdAndUp;
  }
  set drawerProp(value: boolean) {
    if (this.$vuetify.breakpoint.mdAndUp) {
      this.drawer = false;
    } else this.drawer = value;
  }
  onResize(setdrawer: boolean) {
    this.appBarHeight = document.getElementById("appBarElm").clientHeight;
    if (this.$vuetify.breakpoint.smAndDown) {
      this.$vuetify.application.right = 0;
      if (setdrawer) this.drawer = false;
    } else {
      if (setdrawer) this.drawer = true;
      this.$vuetify.application.right = 57;
    }
  }
  mounted() {
    this.onResize(false);
  }
  showLogin() {
    this.drawer = false;
    this.enable = true;
  }
  changeTheme() {
    localStorage.setItem("darkTheme", String(this.$vuetify.theme.dark));
  }
  appBarHeight = 64;
}
</script>
<style>
.clickable {
  cursor: pointer;
}
</style>