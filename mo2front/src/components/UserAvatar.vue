<template>
  <v-tooltip bottom content-class="">
    <template v-slot:activator="{ on, attrs }">
      <v-badge color="transparent" avatar overlap bottom>
        <template v-slot:badge>
          <v-hover
            v-model="hov"
            style="z-index: 3; position: relative"
            v-slot="{ hover }"
          >
            <div v-on:click.prevent v-on:click.stop class="unclickable">
              <v-avatar :style="hover && enableEdit ? 'position:absolute' : ''">
                {{ user.settings.status ? user.settings.status : "😀" }}
              </v-avatar>
              <v-emoji-picker
                :style="mobile ? 'margin-left: -120%!important;' : ''"
                :dark="dark"
                class="ml-4"
                v-on:click.prevent
                v-on:click.stop
                v-if="hover && enableEdit"
                @select="change"
              />
            </div>
          </v-hover>
        </template>
        <v-avatar v-bind="attrs" v-on="on" :size="size" color="brown">
          <v-img
            v-if="
              user.settings &&
              user.settings.avatar &&
              user.settings.avatar !== ''
            "
            :src="user.settings.avatar + '~thumb'"
          ></v-img>
          <span v-else class="white--text headline">{{ initials }}</span>
        </v-avatar>
      </v-badge>
    </template>
    <v-card-title>{{ user.name }}</v-card-title>
    <v-divider />
    <v-card-subtitle>{{
      email
        ? user.email.endsWith("@mo2.com")
          ? "anonymous account"
          : user.email
        : "github oauth account"
    }}</v-card-subtitle>
  </v-tooltip>
</template>

<script lang="ts">
import { User } from "@/models";
import { GetInitials, UserRole } from "../utils/index";
import Vue from "vue";
import Component from "vue-class-component";
import { Prop, Watch } from "vue-property-decorator";
import { VEmojiPicker } from "v-emoji-picker";
import { IEmoji } from "node_modules/v-emoji-picker/lib/models/Emoji";
@Component({
  components: {
    VEmojiPicker,
  },
})
export default class Avatar extends Vue {
  @Prop()
  user!: User;
  @Prop()
  size?: number;
  @Prop({ default: false })
  enableEdit!: boolean;
  hov = false;
  get email() {
    return this.user.email.indexOf("@") > 0;
  }
  get isUser() {
    return this.user.roles && this.user.roles.indexOf(UserRole) >= 0;
  }
  get dark() {
    return this.$vuetify.theme.dark;
  }

  get initials(): string {
    try {
      return GetInitials(this.user.name);
    } catch (error) {
      return "A";
    }
  }
  change(emoji: IEmoji) {
    this.hov = false;
    this.$emit("setemoji", emoji);
  }
  get mobile() {
    return this.$vuetify.breakpoint.mobile;
  }
}
</script>


<style>

</style>
