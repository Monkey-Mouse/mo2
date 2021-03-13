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
          @click="$router.push('/')"
          alt="Logo"
          class="shrink mr-2 clickable"
          contain
          src="./assets/logo.png"
          transition="scale-transition"
          width="100"
        />
        <!-- <div class="text-h4">MO2</div> -->
        <!-- <v-img
          alt="Vuetify Name"
          class="shrink mt-1 hidden-sm-and-down"
          contain
          min-width="100"
          src="https://cdn.vuetifyjs.com/images/logos/vuetify-name-dark.png"
          width="100"
        /> -->
      </div>
      <v-spacer />
      <div v-if="$route.path.indexOf('/edit') === 0">
        <v-row>
          <div v-if="autoSaving" class="grey--text ma-2">
            <v-progress-circular
              color="yellow"
              indeterminate
            ></v-progress-circular>
          </div>
          <div v-else-if="autoSaving === null" class="red--text ma-2">
            Auto Save Failed!
          </div>
          <div v-else class="light-green--text ma-2">Saved!</div>
          <v-btn class="ml-10" color="green" outlined @click="publishClick"
            >publish</v-btn
          >
        </v-row>
      </div>
      <v-app-bar-nav-icon
        v-if="!this.$vuetify.breakpoint.mdAndUp"
        @click.stop="
          () => {
            drawer = !drawer;
          }
        "
      ></v-app-bar-nav-icon>
    </v-app-bar>
    <v-navigation-drawer
      right
      fixed
      :permanent="this.$vuetify.breakpoint.mdAndUp"
      :expand-on-hover="this.$vuetify.breakpoint.mdAndUp"
      v-model="drawerProp"
    >
      <v-list-item class="px-2" :style="`height: ${appBarHeight}px`">
        <v-list-item-avatar>
          <avatar :size="40" :user="user" />
        </v-list-item-avatar>

        <v-list-item-title>{{
          isUser ? user.name : "未登录"
        }}</v-list-item-title>

        <v-btn icon v-if="!isUser" @click="showLogin()"> 登录 </v-btn>
        <v-btn icon v-else color="red" @click="logOut"> 登出 </v-btn>
      </v-list-item>

      <v-divider></v-divider>

      <v-list dense>
        <v-list-item
          v-for="(item, n) in items"
          :key="n"
          :to="item.href"
          v-show="item.show"
          :exact="true"
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
      <router-view
        ref="view"
        v-if="userload"
        :autoSaving.sync="autoSaving"
        :user="user"
      />
      <account-modal :enable.sync="enable" :user.sync="userdata" />
      <v-snackbar v-model="snackbar" :timeout="5000">
        {{ "登出成功！" }}

        <template v-slot:action="{ attrs }">
          <v-btn color="blue" text v-bind="attrs" @click="snackbar = false">
            Close
          </v-btn>
        </template>
      </v-snackbar>
      <!-- <v-btn @click="showLogin()">Login</v-btn>
      <account-modal :enable.sync="enable" />
      <v-btn @click="$vuetify.theme.dark = !$vuetify.theme.dark"
        >switch theme</v-btn
      >
      <v-btn @click="sideNavVisible = !sideNavVisible">show side bar</v-btn> -->
    </v-main>
    <v-footer padless>
      <v-card flat tile class="indigo lighten-1 white--text text-center">
        <v-divider></v-divider>

        <v-card-text class="white--text">
          {{ new Date().getFullYear() }} — <strong>MO2</strong>
        </v-card-text>
      </v-card>
      <a class="ml-16" target="blank" href="http://beian.miit.gov.cn/"
        >冀ICP备20007570号-2</a
      >
    </v-footer>
  </v-app>
</template>

<script lang="ts">
import Vue from "vue";
import AccountModal from "./components/AccountModal.vue";
import Vuelidate from "vuelidate";
import Component from "vue-class-component";
import { BlankUser, User } from "./models";
import {
  GetInitials,
  GetUserInfoAsync,
  Logout,
  ReachedBottom,
  UserRole,
} from "./utils";
import Avatar from "./components/UserAvatar.vue";
import { Watch } from "vue-property-decorator";
import "vue2-timeago/dist/vue2-timeago.css";
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
  user: User = BlankUser;
  autoSaving = false;
  snackbar = false;
  userload = false;
  get userdata() {
    return this.user;
  }

  set userdata(v: User) {
    this.user = v;
    this.items[1].show = this.isUser;
    this.items[2].show = this.isUser;
  }

  enable = false;
  items = [
    { title: "Home", icon: "mdi-home-city", href: "/", show: true },
    {
      title: "My Home",
      icon: "mdi-account",
      href: "/account",
      show: this.isUser,
    },
    {
      title: "New Article",
      icon: "mdi-file-document-edit",
      href: "/edit",
      show: this.isUser,
    },
    { title: "About", icon: "mdi-alpha-a-circle", href: "/about", show: true },
  ];
  get isUser() {
    return this.user.roles && this.user.roles.indexOf(UserRole) >= 0;
  }
  @Watch("user")
  userChange() {
    this.items[2].show = this.isUser;
    this.items[1].show = this.isUser;
  }
  get initials(): string {
    return GetInitials(this.user.name);
  }
  logOut() {
    Logout().then(() => {
      GetUserInfoAsync().then((u) => {
        this.user = u;
        this.items[1].show = this.isUser;
        this.items[2].show = this.isUser;
        this.snackbar = true;
      });
    });
  }
  publishClick() {
    (this.$refs["view"] as any).publish();
  }
  created() {
    document.title = "Mo2";
    GetUserInfoAsync().then((u) => {
      this.user = u;
      this.userload = true;
      this.items[1].show = this.isUser;
      this.items[2].show = this.isUser;
    });
    try {
      this.$vuetify.theme.dark = JSON.parse(
        localStorage.getItem("darkTheme")
      ) as boolean;
    } catch (err) {}
    window.addEventListener("resize", () => {
      this.onResize();
      // setTimeout(() => {
      //   this.onResize();
      // }, 500);
    });
    window.addEventListener("scroll", () => {
      if (ReachedBottom()) {
        try {
          (this.$refs["view"] as any).ReachedButtom();
        } catch (error) {}
      }
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
  onResize() {
    this.appBarHeight = document.getElementById("appBarElm").clientHeight;
    if (!this.$vuetify.breakpoint.mdAndUp) {
      this.$vuetify.application.right = 0;
    } else {
      this.$vuetify.application.right = 57;
    }
  }
  mounted() {
    this.onResize();
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
<style lang="scss">
@import "./assets/main.scss";
</style>
<style>
.clickable {
  cursor: pointer;
}
</style>