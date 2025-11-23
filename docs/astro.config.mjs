import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

// https://astro.build/config
export default defineConfig({
  site: 'https://afadesigns.github.io',
  base: '/zshellcheck',
  integrations: [
    starlight({
      title: 'ZShellCheck',
      social: [
        {
          icon: 'github',
          link: 'https://github.com/afadesigns/zshellcheck',
        },
      ],
      sidebar: [
        {
          label: 'Start Here',
          items: [
            { label: 'Introduction', link: '/' },
          ],
        },
        {
          label: 'Guides',
          items: [
            { label: 'Contributing', link: '/guides/contributing/' },
          ],
        },
        {
          label: 'Project Info',
          items: [
            { label: 'Roadmap', link: '/about/roadmap/' },
            { label: 'Code of Conduct', link: '/about/code-of-conduct/' },
          ],
        },
      ],
    }),
  ],
});
