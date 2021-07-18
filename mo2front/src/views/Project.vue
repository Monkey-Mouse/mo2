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
    <v-row>
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
import { BlogBrief, Project, InputProp, User, Option } from "@/models";
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
} from "@/utils";
import Vue from "vue";
import Component from "vue-class-component";
import MO2Dialog from "../components/MO2Dialog.vue";
import {
  DeleteProject,
  GetProject,
  GetProjectArticles,
  GetUserDatas,
  ListProject,
  SearchUser,
} from "../utils/api";
import BlogTimeLineList from "../components/BlogTimeLineList.vue";
import BlogSkeleton from "../components/BlogTimeLineSkeleton.vue";
import { required } from "vuelidate/lib/validators";
import { Prop, Watch } from "vue-property-decorator";
const lazySearcher = new LazyExecutor();
const dic: { [key: string]: Option } = {};
@Component({
  components: {
    BlogTimeLineList,
    BlogSkeleton,
    MO2Dialog,
  },
})
export default class ProjectPage extends Vue implements AutoLoader<BlogBrief> {
  @Prop()
  user: User;
  datalist: BlogBrief[] = [];
  loading = true;
  firstloading = true;
  page = 0;
  pagesize = 10;
  nomore = false;
  proj: Project = {};
  showGroup = false;
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
      this.proj = proj;
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
  created() {
    GetProject(this.$route.params["id"]).then((re) => {
      this.proj = re;
      this.groupProps.name.default = this.proj.Name;
      this.groupProps.description.default = this.proj.Description;
      this.groupProps.tags.default = this.proj.Tags;
      GetUserDatas(re.ManagerIDs).then((managers) => {
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
