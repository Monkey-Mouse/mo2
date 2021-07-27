<template>
  <v-container >
    <MO2Dialog
      :confirm="updateGroup"
      :confirmText="'Update'"
      :title="'Update Group'"
      :inputProps="groupProps"
      :validator="groupValidator"
      :show.sync="showGroup"
    />
    <v-row class=" pt-16">
      <v-col class=" text-center">
        <v-avatar size="80" >
          <v-img src=https://cdn.vuetifyjs.com/images/cards/docks.jpg></v-img>
        </v-avatar>
      </v-col>
    </v-row>
    <v-row>
      <v-col class=" text-center">
        <div class=" text-h4">{{proj.Name}}
          <v-icon v-if="proj.OwnerID==user.id" @click="showGroup = true"
            >mdi-circle-edit-outline</v-icon>
          <v-icon v-if="proj.OwnerID==user.id" @click="deleteProj"
            >mdi-delete-circle-outline</v-icon>
        </div>
      </v-col>
    </v-row>
    <v-row>
      <v-col class=" text-center">
        <div class=" text-body-1">{{proj.Description}}</div>
      </v-col>
    </v-row>
    <v-row class="mb-2">
      <v-col class=" text-center">
        <v-chip
          v-for="tag in proj.Tags"
          :key="tag"
          class=" ma-1"
        >
          {{ tag }}
        </v-chip>
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="12" md="6">
        <v-row class="mb-2">
          <v-col class="text-h4 text-center">Owner</v-col>
        </v-row>
        <v-row class="mb-3">
          <user-item :user="owner"/>
        </v-row>
      </v-col>
      <v-col cols="12" md="6">
        <v-container >
          <v-row>
            <v-col class="text-h4 text-center">Managers</v-col>
          </v-row>
          <v-row v-if="managers.length>0" class="mb-3">
            <v-col v-for="(v,i) in managers" :key="i">
              <user-item :user="v"/>
            </v-col>
          </v-row>
          <v-row v-else class="mb-3">
            <v-col class="text-center">
              No Manager Yet
            </v-col>
          </v-row>
        </v-container>
      </v-col>
    </v-row>
    <v-container >
      <v-row>
        <v-col class="text-h4 text-center">Members</v-col>
      </v-row>
      <v-row v-if="members.length>0" class="mb-3">
        <v-col v-for="(v,i) in members" :key="i">
          <user-item :user="v"/>
        </v-col>
      </v-row>
      <v-row v-else class="mb-3">
        <v-col class="text-center">
          No Member Yet
        </v-col>
      </v-row>
    </v-container>
    <!-- <v-row>
      <v-col class="text-h4 text-center">Articles</v-col>
    </v-row> -->
    <v-row>
      <v-divider/>
    </v-row>
    <v-row>
      <v-container>
        <blog-time-line-list v-if="!firstloading" :blogs="datalist" />
        <blog-skeleton v-if="loading" :num="pagesize" />
      </v-container>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import {
  BlogBrief,
  Project,
  InputProp,
  User,
  Option,
  UserListData,
  BlankUser,
} from "@/models";
import {
  AddMore,
  AutoLoader,
  ElmReachedBottom,
  UpsertProject,
  GetErrorMsg,
  Prompt,
  LazyExecutor,
  BuildOnUserChange,
  UserFilter,
  logOut,
  ShowLogin,
} from "@/utils";
import Vue from "vue";
import Component from "vue-class-component";
import MO2Dialog from "../components/MO2Dialog.vue";
import {
  DeleteProject,
  GetProject,
  GetProjectArticles,
  GetUserData,
  GetUserDatas,
  GetUserInfoAsync,
  JoinProject,
  ListProject,
  SearchUser,
} from "../utils/api";
import BlogTimeLineList from "../components/BlogTimeLineList.vue";
import BlogSkeleton from "../components/BlogTimeLineSkeleton.vue";
import UserItem from "../components/UserItem.vue";
import { required } from "vuelidate/lib/validators";
import { Prop, Watch } from "vue-property-decorator";
const lazySearcher = new LazyExecutor();
const dic: { [key: string]: Option } = {};
@Component({
  components: {
    BlogTimeLineList,
    BlogSkeleton,
    MO2Dialog,
    UserItem,
  },
})
export default class ProjectPage extends Vue implements AutoLoader<BlogBrief> {
  @Prop()
  user: User;
  owner: User = BlankUser;
  datalist: BlogBrief[] = [];
  loading = true;
  firstloading = true;
  page = 0;
  pagesize = 10;
  nomore = false;
  proj: Project = {};
  showGroup = false;
  managers: UserListData[] = [];
  members: UserListData[] = [];
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
    MemberIDs: {},
    ManagerIDs: {},
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
      label: "Tags",
      default: [],
      icon: "mdi-tag-multiple",
      col: 12,
      type: "combo",
      options: ["课程", "娱乐", "互联网", "教育"],
      message: "enter添加自定义tag",
      multiple: true,
    },
    MemberIDs: {
      errorMsg: {},
      label: "Members",
      default: [],
      icon: "mdi-account-group",
      col: 12,
      type: "select",
      options: [],
      multiple: true,
      showAvatar: true,
      onChange: BuildOnUserChange(lazySearcher, dic),
      filter: UserFilter,
    },
    ManagerIDs: {
      errorMsg: {},
      label: "Managers",
      default: [],
      icon: "mdi-account-supervisor",
      col: 12,
      type: "select",
      options: [],
      multiple: true,
      showAvatar: true,
      onChange: BuildOnUserChange(lazySearcher, dic),
      filter: UserFilter,
    },
  };
  async updateGroup(p: Project): Promise<{ err: string; pass: boolean }> {
    try {
      p.ID = this.proj.ID;
      const proj = await UpsertProject(p);
      this.proj = proj.project;
      if (proj.invite) {
        Prompt("邀请已发送", 5000);
      }
      return { err: null, pass: true };
    } catch (error) {
      return { err: GetErrorMsg(error), pass: false };
    }
  }
  async deleteProj() {
    await DeleteProject(this.proj.ID);
    Prompt("删除成功", 5000);
    this.$router.back();
  }
  async joinProj(query: { token: string }) {
    const p = await JoinProject(query);
    if (p.ManagerIDs?.length > this.proj.ManagerIDs?.length) {
      this.managers.push(this.user);
    } else if (p.MemberIDs?.length > this.proj.MemberIDs?.length) {
      this.members.push(this.user);
    }
    return p;
  }
  created() {
    GetProject(this.$route.params["id"]).then((re) => {
      this.proj = re;
      this.groupProps.name.default = this.proj.Name;
      this.groupProps.description.default = this.proj.Description;
      this.groupProps.tags.default = this.proj.Tags;
      GetUserData(this.proj.OwnerID).then((u) => {
        this.owner = u;
      });
      GetUserDatas(re.ManagerIDs).then((managers) => {
        this.managers = managers;
        this.groupProps.ManagerIDs.default = re.ManagerIDs ?? [];
        for (let index = 0; index < managers.length; index++) {
          const u = managers[index];
          dic[u.id] = { text: u.name, value: u.id, avatar: u.settings?.avatar };
        }
        this.groupProps.ManagerIDs.options = managers.map((u) => {
          return { text: u.name, value: u.id, avatar: u.settings?.avatar };
        });
      });
      GetUserDatas(re.MemberIDs).then((members) => {
        this.members = members;
        for (let index = 0; index < members.length; index++) {
          const u = members[index];
          dic[u.id] = { text: u.name, value: u.id, avatar: u.settings?.avatar };
        }
        this.groupProps.MemberIDs.default = re.MemberIDs ?? [];
        this.groupProps.MemberIDs.options = members.map((u) => {
          return { text: u.name, value: u.id, avatar: u.settings?.avatar };
        });
      });
    });
    GetProjectArticles({
      page: this.page++,
      pageSize: this.pagesize,
      id: this.$route.params["id"],
    }).then((val) => {
      AddMore(this, val);
      this.firstloading = false;
    });
    const token = this.$route.query["token"] as string;
    const email = this.$route.query["email"] as string;
    if (token && email) {
      if (email !== this.user.email) {
        if (this.user.roles.indexOf("OrdinaryUser") > -1) {
          logOut().then(() => {
            ShowLogin(email);
          });
        } else ShowLogin(email);
      } else
        this.joinProj({ token: token }).then((p) => {
          this.proj = p;
        });
    }
  }
  @Watch("user")
  userChange() {
    const token = this.$route.query["token"] as string;
    const email = this.$route.query["email"] as string;
    if (token && email && email === this.user.email) {
      this.joinProj({ token: token }).then((p) => {
        this.proj = p;
      });
    }
  }
  public ReachedButtom() {
    ElmReachedBottom(this, ({ page, pageSize }) =>
      GetProjectArticles({
        page: page,
        pageSize: pageSize,
        id: this.$route.params["id"],
      })
    );
  }
}
</script>
<style>
.bordered {
  border: 2px solid #3298dc;
}
</style>
