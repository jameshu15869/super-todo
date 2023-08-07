import MainContainer from "@/components/MainContainer";
import MainNav from "@/components/MainNav";
import Footer from "@/components/Footer/Footer";

export default async function Home() {
  return (
    <>
      <MainNav />
      <MainContainer className="pt-16" />
      <Footer />
    </>
  );
}
