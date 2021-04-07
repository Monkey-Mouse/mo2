<template>
  <v-app>
    <v-app-bar id="appBarElm" color="primary" app>
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
      </div>
      <v-spacer />
      <v-btn
        v-if="this.$route.name !== 'Search Article' && !search"
        @click="search = true"
        icon
        color="secondary"
      >
        <v-icon>mdi-magnify</v-icon>
      </v-btn>
      <v-expand-transition v-if="this.$route.name !== 'Search Article'">
        <v-autocomplete
          color="secondary"
          v-if="search"
          autofocus
          @blur="search = false"
          class="pa-1"
          style="max-width: 300px"
          label="Search"
          clearable
          :search-input.sync="searchString"
          hide-details="auto"
          @keydown="keyDown"
          :items="prompts"
          no-filter
          item-text="title"
          item-value="id"
        >
          <template v-slot:item="{ item }">
            <v-card
              v-on:click.prevent
              v-on:click.stop
              class="ma-3"
              style="max-width: 300px; width: 100%"
              @click="
                $router.push('/article/' + item.id);
                searchString = '';
                search = false;
              "
            >
              <div class="d-flex flex-no-wrap justify-space-between">
                <div>
                  <v-card-title
                    class="headline"
                    v-html="item.title"
                  ></v-card-title>

                  <v-card-subtitle
                    class="ellipsis-1"
                    v-html="item.description"
                  ></v-card-subtitle>
                </div>

                <v-avatar class="ma-3" size="64" tile>
                  <v-img :src="item.cover"></v-img>
                </v-avatar>
              </div>
            </v-card>
          </template>
        </v-autocomplete>
      </v-expand-transition>
      <v-btn
        v-if="!isUser && !search"
        color="success"
        outlined
        @click="showLogin()"
      >
        LOGIN
      </v-btn>
      <div v-if="$route.path.indexOf('/edit') === 0">
        <v-row>
          <div v-if="autoSaving" class="grey--text ma-2">
            <v-progress-circular
              color="yellow"
              indeterminate
            ></v-progress-circular>
          </div>
          <div v-else-if="autoSaving === null" class="error--text ma-2">
            Failed!
          </div>
          <div v-else class="success--text ma-2">Saved!</div>
          <v-btn class="ml-1" color="success" outlined @click="publishClick"
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
      <v-progress-linear
        :active="pos !== 0"
        :value="pos"
        absolute
        bottom
        color="black white"
      ></v-progress-linear>
    </v-app-bar>
    <v-navigation-drawer
      right
      fixed
      :permanent="this.$vuetify.breakpoint.mdAndUp"
      :expand-on-hover="this.$vuetify.breakpoint.mdAndUp"
      v-model="drawerProp"
    >
      <v-list-item class="px-2" :style="`height: ${appBarHeight}px`">
        <v-badge
          :value="notificationNum !== 0"
          color="red"
          :content="notificationNum"
          offset-x="15"
          offset-y="10"
        >
          <div
            :class="notificationNum !== 0 ? 'clickable' : ''"
            @click="
              () => {
                if (notificationNum !== 0) {
                  $router.push('/notifications');
                  notificationNum = 0;
                }
              }
            "
          >
            <avatar :size="40" :user="user" />
          </div>
        </v-badge>

        <v-list-item-title class="ml-7">{{
          isUser ? user.name : "未登录"
        }}</v-list-item-title>

        <v-btn
          icon
          v-if="isUser"
          @click="
            $router.push('/notifications');
            notificationNum = 0;
          "
        >
          <v-icon>mdi-email</v-icon>
        </v-btn>
        <v-btn icon v-if="!isUser" @click="showLogin()"> 登录 </v-btn>
        <v-btn icon v-else color="error" @click="logOut"> 登出 </v-btn>
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
        <v-list-item to="/settings">
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
      <v-snackbar v-model="refresh" timeout="-1">
        发现新版本，请刷新

        <template v-slot:action="{ attrs }">
          <v-btn color="accent" text v-bind="attrs" @click="reload">
            Refresh
          </v-btn>
          <v-btn color="pink" text v-bind="attrs" @click="refresh = false">
            Close
          </v-btn>
        </template>
      </v-snackbar>
      <v-snackbar v-model="prompt" :timeout="ptimeout">
        {{ pmsg }}

        <template v-slot:action="{ attrs }">
          <v-btn color="pink" text v-bind="attrs" @click="prompt = false">
            Close
          </v-btn>
        </template>
      </v-snackbar>
    </v-main>
    <v-footer id="footer" padless>
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
import { BlankUser, BlogBrief, User } from "./models";
import {
  GetArticles,
  GetInitials,
  GetNotificationNums,
  GetUserInfoAsync,
  LazyExecutor,
  Logout,
  ReachedBottom,
  SetApp,
  SetTheme,
  SetThemeColors,
  ShowRefresh,
  SlowExecutor,
  UserRole,
} from "./utils";
import Avatar from "./components/UserAvatar.vue";
import { Watch } from "vue-property-decorator";
import "vue2-timeago/dist/vue2-timeago.css";
import { VuetifyThemeVariant } from "vuetify/types/services/theme";
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
  search = false;
  searchString = "";
  refresh = false;
  prompts: BlogBrief[] = [];
  searchLoader: LazyExecutor = new LazyExecutor(null, 200);
  pos = 0;

  prompt = false;
  pmsg = "";
  ptimeout = 5000;
  Prompt(msg: string, timeout: number) {
    this.pmsg = msg;
    this.ptimeout = timeout;
    this.prompt = true;
  }

  keyDown(event: KeyboardEvent) {
    // Number 13 is the "Enter" key on the keyboard
    if (event.key === "Enter") {
      // Cancel the default action, if needed
      event.preventDefault();
      // Trigger the button element with a click
      this.search = false;
      this.$router
        .push("/search?q=" + this.searchString)
        .then(() => (this.searchString = ""));
    } else
      this.searchLoader.Execute(() =>
        GetArticles({
          page: 0,
          pageSize: 5,
          draft: false,
          search: this.searchString,
        }).then((bs) => {
          this.prompts = bs;
        })
      );
  }
  reload() {
    this.refresh = false;
    window.location.reload();
  }
  get userdata() {
    return this.user;
  }
  set userdata(v: User) {
    this.user = v;
    this.items[1].show = this.isUser;
    this.items[2].show = this.isUser;
  }
  enable = false;
  notificationNum = 0;
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
    try {
      if (this.isUser) {
        GetNotificationNums().then((d) => {
          this.notificationNum = d.num;
        });
      }
      if (this.user.settings && this.user.settings.perferDark) {
        SetTheme(JSON.parse(this.user.settings.perferDark) as boolean, this);
        if (this.user.settings.themes) {
          const theme = JSON.parse(this.user.settings.themes) as {
            light: VuetifyThemeVariant;
            dark: VuetifyThemeVariant;
          };
          if (!theme) {
            return;
          }
          SetThemeColors(this, theme);
          SetTheme(
            JSON.parse(this.user.settings.perferDark) as boolean,
            this,
            theme
          );
        }
      }
    } catch (error) {
      console.error(error);
    }
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
    window.addEventListener("scroll", () => {
      const h = document.documentElement,
        b = document.body,
        st = "scrollTop",
        sh = "scrollHeight";
      this.pos = ((h[st] || b[st]) / ((h[sh] || b[sh]) - h.clientHeight)) * 100;
    });
    SetApp(this);
    document.title = "Mo2";
    this.$router.afterEach(() => {
      this.pos = 0;
    });
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
    } catch (err) {
      console.error(err);
    }
    try {
      const themes = JSON.parse(localStorage.getItem("themes")) as {
        light: VuetifyThemeVariant;
        dark: VuetifyThemeVariant;
      };
      if (themes) {
        SetThemeColors(this, themes);
      }
    } catch (err) {
      console.error(err);
    }
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
        } catch (error) {
          console.error(error);
        }
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
    SetTheme(
      this.$vuetify.theme.dark,
      this,
      this.$vuetify.theme.themes,
      this.user
    );
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
.anchor {
  padding-top: 70px !important;
  margin-top: -70px !important;
}
</style>