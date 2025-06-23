<script setup lang="ts">
import { Info, User } from 'lucide-vue-next'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { BaseInfoForm, FormActions, FormCard, PageHeader } from '@/components'
import { useSecretItemForm } from '@/composables/useSecretItemForm'
import { SECRET_ITEM_TYPE } from '@/constants'
import PasswordDataForm from '@/views/password/PasswordDataForm.vue'

const router = useRouter()
const route = useRoute()

// 是否为编辑模式
const isEdit = computed(() => !!route.params.id)

const {
  formData,
  loading,
  handleSubmit,
} = useSecretItemForm(SECRET_ITEM_TYPE.Password)
</script>

<template>
  <div class="space-y-6">
    <PageHeader
      :title="isEdit ? '编辑账号密码' : '新建账号密码'" :description="isEdit ? '修改现有的账号密码信息' : '创建一个新的账号密码'"
      @back="router.back()"
    />

    <form @submit="handleSubmit">
      <div class="grid gap-6 lg:grid-cols-2">
        <!-- 基础信息 -->
        <FormCard title="基础信息" description="填写账号密码的基本信息" :icon="Info">
          <BaseInfoForm v-model="formData" :type="SECRET_ITEM_TYPE.Password" />
        </FormCard>

        <!-- 账号密码信息 -->
        <FormCard title="账号信息" description="配置账号密码的具体信息" :icon="User">
          <PasswordDataForm />
        </FormCard>
      </div>
      <FormActions :loading="loading" :is-edit="isEdit" @cancel="router.back()" />
    </form>
  </div>
</template>
