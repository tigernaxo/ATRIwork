#include <string>
#include <array>
#include <fstream>
#include <iostream>
#include <iterator>
#include <algorithm>
#include <vector>
#include <utility>
#include <tuple>
#include <functional>

using namespace std;
using namespace std::placeholders;

pair<unsigned int, string> countGetRef(const string f);
void setMask(string file, string& refSeq,unsigned int&  minLen, bool* nMask, bool* snpMask);
void outputMaskFasta(string file, const unsigned int minLen, bool* totalMask);
void printHelp();

int main(int argc, char* argv[])
{ 
  if(argc == 1)
  {
    printHelp();
    return 0;
  }

  // get file list
  vector<string> fileList;
  for(size_t i = 1; i < argc; ++i)
  {
    fileList.push_back(argv[i]);
  }

  // Determine mask length then get ref seq.
  string refSeq;
  unsigned int len, minLen;
  std::tie(len, refSeq) = countGetRef(fileList[0]);
  minLen = len;

  // create 3 bool vector, length = len
  bool* nMask = new bool[len];
  bool* snpMask = new bool[len];
  bool* totalMask = new bool[len];
  for(size_t i=0; i != len; ++i)
  {
    nMask[i] = true;
    snpMask[i] = false;
  }

  // set mask for each file
  //for_each(fileList.begin(), fileList.end(), setMaskLambda);
  for(auto it=fileList.begin(); it !=fileList.end(); ++it)
  {
    setMask(*it , refSeq, minLen, nMask, snpMask);
  }

  // set Total Mask
  for(size_t i = 0; i != minLen; ++i)
  {
    if(snpMask[i] && nMask[i])
    {
      totalMask[i] = true;
    }else{
      totalMask[i] = false;
    }
  }

  // output masked fasta
  for(auto it=fileList.begin(); it !=fileList.end(); ++it)
  {
    outputMaskFasta(*it, minLen, totalMask);
  }
  
  return 0;
}

void printHelp()
{
    cout << endl << endl;
    cout << "====================snpAligner version 1.0.1======================" << endl << endl;
    cout << "snpAligner do folowing things:" << endl;
    cout << "    1.find positions \"do not have n/N\" in all input file." << endl;
    cout << "    2.find positions \"have snp\" in all input file." << endl;
    cout << "    3.set output sites = \"don't have n/N\" and \"have snp\" in all input fasta."<< endl;
    cout << "    4.Filter all input fasta to stdout depend on output sites, 1 fasta 2 lines (id ,sequence). "<< endl;
    cout << "    5.snpAligner use input file name (without extension) as output fasta id, not input fasta id." << endl;
    cout << endl;
    cout << "usage:" << endl;
    cout << "    snpAilgner [input.fastas...] > output.fasta" << endl;
    cout << endl;
    cout << "example:" << endl;
    cout << "    snpAilgner fasta1.fasta fasta2.fasta fasta3.fasta > output.fasta" << endl;
    cout << endl;
    cout << "output.fasta:" << endl;
    cout << ">fasta1" << endl;
    cout << "atagatagaacg..." << endl;
    cout << ">fasta2" << endl;
    cout << "agctagctgtgg..." << endl;
    cout << ">fasta3" << endl;
    cout << "ctcgctccatct..." << endl;
    cout << "    " << endl;
    cout << "    Have a fun time     Yu-Cheng,Chen" << endl;
    cout << endl << endl;
}
pair<unsigned int, string> countGetRef(const string f)
{
  ifstream ifs(f);
  unsigned int count(0);
  string s;
  char c;

  ifs.seekg(0, ifs.end);
  string ref(ifs.tellg(), 'n');
  ifs.seekg(0, ifs.beg);

  getline(ifs, s);
  while(ifs.get(c))
  {
    if(c != 10){
      ref[count] = c;
      count++;
    }
  }
  ifs.close();

  ref.resize(count);
  return make_pair( count, ref);
}
void setMask(string file, string& refSeq,unsigned int&  minLen, bool* nMask, bool* snpMask)
{
  ifstream ifs(file);
  unsigned int count(0);
  string s;
  char c;

  getline(ifs, s);
  while(ifs.get(c))
  {
    if(c != 10)
    {
      if(nMask[count])
      {
        if(c == 78 || c == 110)
        {
          nMask[count] = false;
        }else if(c >= 65 && c <= 90)
        {
          if(c != refSeq[count] && c != refSeq[count] - 32)
          {
            snpMask[count] = true;
          }
        }else if(c >= 97 && c <= 122)
        {
          if(c != refSeq[count] && c != refSeq[count] + 32)
          {
            snpMask[count] = true;
          }
        }else{
          if(c != refSeq[count])
            snpMask[count] = true;
        }
      }
      if(++count >= minLen)
        break;
    }
  }
  ifs.close();
}
void outputMaskFasta(string file, const unsigned int minLen, bool* totalMask)
{
  string::size_type start, length;
  start = file.rfind("/") + 1;
  length = file.rfind(".") - start;
  cout << ">" << file.substr(start, length) << "\n";
  ifstream ifs(file);
  unsigned int count(0);
  string s;
  char c;

  getline(ifs, s);
  while(ifs.get(c))
  {
    if(c != 10)
    {
      if(totalMask[count])
        cout << c;
        ++count;
    }
    if(count >= minLen)
      break;
  }
  cout << '\n';
  ifs.close();
}
