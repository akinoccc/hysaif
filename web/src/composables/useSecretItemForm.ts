import type { SecretItem } from '@/api/types'
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import * as z from 'zod'
import { secretItemAPI } from '@/api/secret'
import { SECRET_ITEM_TYPE } from '@/constants'

export function useSecretItemForm(itemType: string) {
  const router = useRouter()
  const route = useRoute()
  const loading = ref(false)
  const isEdit = computed(() => !!route.params.id)

  // 根据不同的密钥类型定义不同的表单验证模式
  let dataSchema: any

  // 根据不同的密钥类型设置不同的初始值
  let initialData: any

  switch (itemType) {
    case SECRET_ITEM_TYPE.ApiKey:
      dataSchema = z.object({
        api_key: z.string().min(1, '请输入 API 密钥'),
        api_secret: z.string().optional(),
        endpoint: z.string().optional(),
        rate_limit: z.string().optional(),
        scope: z.string().optional(),
        notes: z.string().optional(),
      })

      initialData = {
        api_key: '',
        api_secret: '',
        endpoint: '',
        rate_limit: '',
        scope: '',
        notes: '',
      }
      break
    case SECRET_ITEM_TYPE.AccessKey:
      dataSchema = z.object({
        access_key: z.string().min(1, '请输入 Access Key'),
        secret_key: z.string().min(1, '请输入 Secret Key'),
        region: z.string().optional(),
        notes: z.string().optional(),
      })
      initialData = {
        access_key: '',
        secret_key: '',
        region: '',
        notes: '',
      }
      break
    case SECRET_ITEM_TYPE.Password:
      dataSchema = z.object({
        username: z.string().min(1, '请输入用户名'),
        password: z.string().min(1, '请输入密码'),
        address: z.string().optional(),
        notes: z.string().optional(),
      })
      initialData = {
        username: '',
        password: '',
        address: '',
        notes: '',
      }
      break
    case SECRET_ITEM_TYPE.SshKey:
      dataSchema = z.object({
        private_key: z.string().min(1, '请输入私钥'),
        public_key: z.string().optional(),
        passphrase: z.string().optional(),
        key_type: z.string().optional(),
        key_size: z.string().optional(),
        fingerprint: z.string().optional(),
        notes: z.string().optional(),
      })
      initialData = {
        private_key: '',
        public_key: '',
        passphrase: '',
        key_type: '',
        key_size: '',
        fingerprint: '',
        notes: '',
      }
      break
    case SECRET_ITEM_TYPE.Token:
      dataSchema = z.object({
        token: z.string().min(1, '请输入令牌'),
        token_type: z.string().optional(),
        refresh_token: z.string().optional(),
        notes: z.string().optional(),
      })
      initialData = {
        token: '',
        token_type: '',
        refresh_token: '',
        notes: '',
      }
      break
    case SECRET_ITEM_TYPE.Custom:
      dataSchema = z.object({
        custom_data: z.array(z.object({
          key: z.string().min(1, '请输入键'),
          value: z.string().min(1, '请输入值'),
        })).min(1, '请添加自定义数据'),
        notes: z.string().optional(),
      })
      initialData = {
        custom_data: [],
        notes: '',
      }
      break
    default:
      dataSchema = z.object({
        notes: z.string().optional(),
      })
      initialData = {
        notes: '',
      }
      break
  }

  // 表单验证模式
  const formSchema = toTypedSchema(
    z.object({
      name: z.string().min(1, '请输入名称'),
      description: z.string().optional(),
      category: z.string().optional(),
      tags: z.array(z.string()).optional(),
      expires_at: z.number().optional(),
      data: dataSchema,
    }),
  )

  // 表单实例
  const form = useForm({
    validationSchema: formSchema,
    initialValues: {
      name: '',
      description: '',
      category: '',
      tags: [],
      expires_at: undefined,
      data: initialData,
    },
  })

  // 基础信息数据绑定
  const formData = computed({
    get: () => ({
      name: form.values.name,
      description: form.values.description,
      category: form.values.category,
      tags: form.values.tags,
      expires_at: form.values.expires_at,
    }),
    set: (value) => {
      form.setFieldValue('name', value.name)
      form.setFieldValue('description', value.description)
      form.setFieldValue('category', value.category)
      form.setFieldValue('tags', value.tags)
      form.setFieldValue('expires_at', value.expires_at)
    },
  })

  // 加载现有数据（编辑模式）
  const loadItem = async () => {
    if (!isEdit.value)
      return

    try {
      loading.value = true
      const response = await secretItemAPI.getItem(route.params.id as string)
      const item = response as SecretItem

      // 根据不同的密钥类型准备不同的数据对象
      let itemData: any = {}

      if (itemType === SECRET_ITEM_TYPE.ApiKey) {
        itemData = {
          api_key: item.data?.api_key,
          api_secret: item.data?.api_secret,
          notes: item.data?.notes,
        }
      }
      else if (itemType === SECRET_ITEM_TYPE.AccessKey) {
        itemData = {
          access_key: item.data?.access_key,
          secret_key: item.data?.secret_key,
          region: item.data?.region,
          notes: item.data?.notes,
        }
      }
      else if (itemType === SECRET_ITEM_TYPE.Password) {
        itemData = {
          username: item.data?.username,
          password: item.data?.password,
          notes: item.data?.notes,
        }
      }
      else if (itemType === SECRET_ITEM_TYPE.SshKey) {
        itemData = {
          private_key: item.data?.private_key,
          public_key: item.data?.public_key,
          passphrase: item.data?.passphrase,
          key_type: item.data?.key_type,
          key_size: item.data?.key_size,
          fingerprint: item.data?.fingerprint,
          notes: item.data?.notes,
        }
      }
      else if (itemType === SECRET_ITEM_TYPE.Token) {
        itemData = {
          token: item.data?.token,
          token_type: item.data?.token_type,
          refresh_token: item.data?.refresh_token,
          notes: item.data?.notes,
        }
      }
      else if (itemType === SECRET_ITEM_TYPE.Custom) {
        itemData = {
          custom_data: Array.isArray(item.data?.custom_data)
            ? item.data.custom_data
            : [{ key: '', value: '' }],
          notes: item.data?.notes,
        }
      }
      else {
        itemData = {
          notes: item.data?.notes,
        }
      }

      // 填充表单数据
      form.setValues({
        name: item.name,
        description: item.description,
        category: item.category,
        tags: item.tags || [],
        expires_at: item.expires_at,
        data: itemData,
      })
    }
    catch (error) {
      console.error('加载数据失败:', error)
    }
    finally {
      loading.value = false
    }
  }

  // 提交表单
  const handleSubmit = form.handleSubmit(async (values) => {
    try {
      loading.value = true

      const payload = {
        name: values.name,
        description: values.description,
        type: itemType,
        category: values.category,
        tags: values.tags,
        expires_at: values.expires_at,
        data: values.data,
      }

      let itemId = route.params.id as string

      if (isEdit.value) {
        await secretItemAPI.updateItem(route.params.id as string, payload)
        toast.success('保存成功')
      }
      else {
        const { id } = await secretItemAPI.createItem(payload)
        itemId = id
        toast.success('创建成功')
      }
      router.push(`/${itemType}/${itemId}`)
    }
    catch (error) {
      console.error('保存失败:', error)
      toast.error('保存失败')
    }
    finally {
      loading.value = false
      loading.value = false
    }
  })

  // 自动加载数据
  onMounted(() => {
    if (isEdit.value) {
      loadItem()
    }
  })

  return {
    form,
    loading,
    isEdit,
    formData,
    loadItem,
    handleSubmit,
    router,
  }
}
