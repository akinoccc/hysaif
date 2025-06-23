import 'vue-router'

declare module 'vue-router' {
  interface RouteMeta {
    menu?: {
      showInMenu?: boolean
      icon?: FunctionalComponent<LucideProps, {}, any, {}>
      title: string
      order?: number
    }
  }
}
