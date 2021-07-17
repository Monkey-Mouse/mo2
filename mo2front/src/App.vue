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
        style="margin-right: 10px"
        @click="showLogin()"
      >
        LOGIN
      </v-btn>
      <div
        v-if="
          $route.path.indexOf('/edit') === 0 &&
          !search &&
          autoSaving !== 'notme'
        "
      >
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
          <v-btn
            class="ml-1 mt-1"
            color="success"
            small
            outlined
            @click="publishClick"
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
      <!-- <v-progress-linear
        :active="true"
        :value="10"
        absolute
        bottom
        color="black white"
      ></v-progress-linear> -->
      <div class="progress-container">
        <div class="progress-bar" id="myBar"></div>
      </div>
    </v-app-bar>
    <v-navigation-drawer
      right
      fixed
      :permanent="this.$vuetify.breakpoint.mdAndUp"
      :expand-on-hover="this.$vuetify.breakpoint.mdAndUp"
      v-model="drawerProp"
    >
      <v-list-item
        class="px-2"
        :style="`height: ${appBarHeight === 0 ? 64 : appBarHeight}px`"
      >
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
          @click="navClick(item)"
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
        <v-list-item v-if="showInstall" @click="install">
          <v-list-item-icon>
            <v-icon>mdi-download-circle</v-icon>
          </v-list-item-icon>

          <v-list-item-content>
            <v-list-item-title>Install</v-list-item-title>
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
      <MO2Dialog
        :confirm="newGroup"
        :confirmText="'Create'"
        :title="'New Group'"
        :inputProps="groupProps"
        :validator="groupValidator"
        :show.sync="showGroup"
      />
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
    <pwa-install></pwa-install>
  </v-app>
</template>

<script lang="ts">
import Vue from "vue";
import AccountModal from "./components/AccountModal.vue";
import MO2Dialog from "./components/MO2Dialog.vue";
import Vuelidate from "vuelidate";
import Component from "vue-class-component";
import { BlankUser, BlogBrief, InputProp, Project, User } from "./models";
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
  UserRole,
  NewGroup,
  UpsertProject,
  GetErrorMsg,
} from "./utils";
import Avatar from "./components/UserAvatar.vue";
import { Watch } from "vue-property-decorator";
import "vue2-timeago/dist/vue2-timeago.css";
import { VuetifyThemeVariant } from "vuetify/types/services/theme";
import "@pwabuilder/pwainstall";
import { required } from "vuelidate/lib/validators";
// import "bulma/bulma.sass";
Vue.use(Vuelidate);

@Component({
  components: {
    AccountModal,
    Avatar,
    MO2Dialog,
  },
})
export default class App extends Vue {
  drawer = false;
  appBarHeight = 64;
  user: User = BlankUser;
  autoSaving: boolean | "notme" = "notme";
  snackbar = false;
  userload = false;
  search = false;
  searchString = "";
  refresh = false;
  prompts: BlogBrief[] = [];
  showGroup = false;
  searchLoader: LazyExecutor = new LazyExecutor(null, 200);

  prompt = false;
  pmsg = "";
  ptimeout = 5000;
  showInstall = false;

  groupValidator = {
    name: {
      required: required,
    },
    description: {
      required: required,
    },
    tags: {
      required: required,
    },
  };
  groupProps: { [name: string]: InputProp } = {
    name: {
      errorMsg: {
        required: "组名不可为空",
      },
      label: "Name",
      default: "",
      icon: "mdi-rename-box",
      col: 12,
      type: "text",
    },
    description: {
      errorMsg: {
        required: "组描述不可为空",
      },
      label: "Description",
      default: "",
      icon: "mdi-text",
      col: 12,
      type: "textarea",
    },
    tags: {
      errorMsg: {
        required: "标签不可为空",
      },
      label: "Description",
      default: [],
      icon: "mdi-text",
      col: 12,
      type: "combo",
      options: ["课程", "娱乐", "互联网", "教育"],
      message: "enter添加自定义tag",
      multiple: true,
    },
  };
  async newGroup(p: Project): Promise<{ err: string; pass: boolean }> {
    try {
      const proj = await UpsertProject(p);
      this.$router.push(`/project/${proj.ID}`);
      return { err: null, pass: true };
    } catch (error) {
      return { err: GetErrorMsg(error), pass: false };
    }
  }
  navClick(item) {
    if (item.click) item.click();
    else this.$router.push(item.href);
  }
  get installComponent(): any {
    return document?.querySelector("pwa-install");
  }
  install() {
    this.installComponent.manifestpath = "/manifest.json";
    this.installComponent.openPrompt();
  }
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
      show: true,
    },
    {
      title: "New Article",
      icon: "mdi-file-document-edit",
      href: "/edit",
      show: true,
    },
    {
      title: "New Group",
      icon: "mdi-account-group",
      click: () => {
        NewGroup();
      },
      show: true,
    },
    {
      title: "Recycle Bin",
      icon: "mdi-delete",
      href: "/recycle",
      show: true,
    },
    { title: "About", icon: "mdi-alpha-a-circle", href: "/about", show: true },
  ];
  get isUser() {
    return this.user.roles && this.user.roles.indexOf(UserRole) >= 0;
  }
  @Watch("user")
  userChange() {
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
        this.snackbar = true;
      });
    });
  }
  publishClick() {
    (this.$refs["view"] as any).publish();
  }
  created() {
    SetApp(this);
    document.title = "Mo2";
    GetUserInfoAsync().then((u) => {
      this.user = u;
      this.userload = true;
      (window as any).loading_screen.finish();
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
      if (
        ReachedBottom() &&
        (this.$refs["view"] as any) &&
        (this.$refs["view"] as any).ReachedButtom
      ) {
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
    const bar = document.getElementById("myBar");
    window.addEventListener("scroll", () => {
      const winScroll =
        document.body.scrollTop || document.documentElement.scrollTop;
      const height =
        document.documentElement.scrollHeight -
        document.documentElement.clientHeight;
      const scrolled = (winScroll / height) * 100;
      if (scrolled < 1) {
        bar.parentElement.style.display = "none";
        return;
      } else {
        bar.parentElement.style.display = "block";
      }
      bar.style.width = scrolled + "%";
    });
    setTimeout(() => {
      if (!localStorage.getItem("install")) this.install();
      localStorage.setItem("install", "prompted");
      this.showInstall = !this.installComponent.getInstalledStatus();
    }, 2000);
  }
  showLogin() {
    if (this.isUser) {
      return;
    }
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
[anchor] {
  padding-top: 70px !important;
  margin-top: -70px !important;
}
</style>