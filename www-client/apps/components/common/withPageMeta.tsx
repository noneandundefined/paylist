import PageMeta from '@/components/PageMeta/PageMeta';

export const withPageMeta = (Page: () => React.JSX.Element, titleKey: string, descriptionKey?: string) => {
	const WrappedPage = () => (
		<>
			<PageMeta titleKey={titleKey} descriptionKey={descriptionKey} />
			<Page />
		</>
	);

	return WrappedPage;
};
